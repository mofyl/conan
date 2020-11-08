package bbr

import (
	"context"
	"errors"
	"conan/rolling"
	cpustate "conan/sys/cpu"
	"math"
	"sync/atomic"
	"time"
)

var (
	cpu         int64
	decay       = 0.95 // 这里这个值是看自己项目进行调整的
	initTime    = time.Now()
	defaultConf = &Config{
		Window:       time.Second * 10,
		WinBucket:    100,
		CPUThreshold: 800,
	}
)

func Init() {
	go cpuProc()
}

type cpuGetter func() int64

func cpuProc() {
	ticker := time.Ticker{}
	for range ticker.C {
		s := cpustate.State{}
		cpustate.ReadState(&s)
		prevCpu := atomic.LoadInt64(&cpu)
		// EWMA
		curCpu := int64(float64(prevCpu)*decay + (1.0-decay)*float64(s.Usage))
		atomic.StoreInt64(&cpu, curCpu)
	}
}

type Config struct {
	Enable       bool
	Window       time.Duration
	WinBucket    int
	CPUThreshold int64
}

type Stat struct {
	CPU       int64
	InFlight  int64
	MaxFlight int64
	MinRt     int64
	MaxPass   int64
}

type BBR struct {
	cpu             cpuGetter
	passStat        rolling.RollingCounter // 该窗口用于统计通过的 请求数量。也就是qps
	rtStat          rolling.RollingCounter // 该窗口用于统计 每次请求的rtt
	inflight        int64                  // 当前正在请求的数量
	winBucketPreSec int64
	config          *Config
	prevDrop        atomic.Value // 上次丢弃的时间 。这里的时间指 time.Since(上次计算时间) 这里保存的是一个time.Duration
	// prePropHit      int32 // 不知道有什么用 感觉可以不要
	rawMaxPASS int64 // 最大通过的请求数量
	rawMinRt   int64 // 最小的请求处理时间
}

func (b *BBR) maxPass() int64 {

	rawMaxPass := atomic.LoadInt64(&b.rawMaxPASS)

	// 这里若是 取值的时间太短 每次计算都是一样的数值 直接取上次的结果即可
	if rawMaxPass > 0 && b.passStat.Timespan() < 1 {
		return rawMaxPass
	}
	rawMaxPass = int64(b.passStat.Reduce(func(iterator rolling.Iterator) float64 {
		var result = 1.0

		for i := 1; iterator.Next() && i < b.config.WinBucket; i++ {
			bucket := iterator.Bucket()

			count := 0.0

			for _, p := range bucket.Point {
				count += p
			}
			result = math.Max(result, count)
		}
		return result
	}))

	if rawMaxPass == 0 {
		rawMaxPass = 1
	}

	atomic.StoreInt64(&b.rawMaxPASS, rawMaxPass)
	return rawMaxPass
}

func (b *BBR) minRt() int64 {
	rawMinRt := atomic.LoadInt64(&b.rawMinRt)
	if rawMinRt > 0 && b.rtStat.Timespan() < 1 {
		return rawMinRt
	}

	rawMinRt = int64(b.rtStat.Reduce(func(iterator rolling.Iterator) float64 {
		var result = math.MaxFloat64

		for i := 1; iterator.Next() && i < b.config.WinBucket; i++ {
			bucket := iterator.Bucket()
			if bucket.Count == 0 {
				continue
			}
			total := 0.0

			for _, p := range bucket.Point {
				total += p
			}
			avg := total / float64(bucket.Count)
			result = math.Min(avg, result)

		}
		return result
	}))

	if rawMinRt <= 0 {
		rawMinRt = 1
	}
	atomic.StoreInt64(&b.rawMinRt, rawMinRt)
	return rawMinRt
}

// maxFlight 公式 = maxQps * minRtt  可以加盐值来调整
// 该值一般的范围为 cpuCores * 2.5
func (b *BBR) maxFlight() int64 {
	// 这里向上取整
	return int64(math.Floor(float64(b.maxPass()*b.minRt()*b.winBucketPreSec)/1000.0 + 0.5))
}

func (b *BBR) shouldDrop() bool {
	// 这里表示 cpu使用率 小于 当前设置的 阈值
	if b.cpu() < b.config.CPUThreshold {
		prevDrop, _ := b.prevDrop.Load().(time.Duration)
		// 这里表示从来没有丢弃过
		if prevDrop == 0 {
			return false
		}

		// 这里表示 两次 限流的时间差 不到1s
		// 这里就要防止 瞬时流量的 激增了
		if time.Since(initTime)-prevDrop <= time.Second {
			// 这里去估算 是否要丢弃
			inFlight := atomic.LoadInt64(&b.inflight)
			return inFlight > 1 && inFlight > b.maxFlight()
		}
		// 这里表示不是激增的情况
		b.prevDrop.Store(time.Duration(0))
		return false
	}
	// 这里就表示了 当前的cpu 使用率超过了阈值 可能超负载了 可能需要Drop了
	inFlight := atomic.LoadInt64(&b.inflight)
	drop := inFlight > b.maxFlight()
	if drop {
		b.prevDrop.Store(time.Since(initTime))
	}
	//if drop {
	//	prevDrop, _ := b.prevDrop.Load().(time.Duration)
	//	if prevDrop != 0 {
	//		return drop
	//	}
	//  b.prevDrop.Store(time.Since(initTime))
	//}
	return drop
}

func (b *BBR) Allow(ctx context.Context) (func(info DoneInfo), error) {
	if b.shouldDrop() {
		return nil, errors.New("should drop")
	}
	atomic.AddInt64(&b.inflight, 1)
	startTime := time.Since(initTime)

	return func(info DoneInfo) {
		// 计算本次的rt
		rt := float64((time.Since(initTime) - startTime) / time.Millisecond)
		b.rtStat.Add(rt)
		atomic.AddInt64(&b.inflight, -1)
		switch info.Op {
		case Success:
			b.passStat.Add(1)
		default:
			return
		}
	}, nil
}

func newLimiter(conf *Config) Limiter {
	if conf == nil {
		conf = defaultConf
	}

	bucketDuration := conf.Window / time.Duration(conf.WinBucket)

	passStat := rolling.NewRollingCounter(rolling.RollingCounterOpts{
		Size:           conf.WinBucket,
		BucketDuration: bucketDuration,
	})
	rtStat := rolling.NewRollingCounter(rolling.RollingCounterOpts{
		Size:           conf.WinBucket,
		BucketDuration: bucketDuration,
	})
	cpuFunc := func() int64 {
		return atomic.LoadInt64(&cpu)
	}

	limiter := &BBR{
		cpu:             cpuFunc,
		passStat:        passStat,
		rtStat:          rtStat,
		winBucketPreSec: int64(time.Second) / int64(bucketDuration),
		config:          conf,
	}
	return limiter
}

func (b *BBR) Stat() Stat {
	return Stat{
		CPU:       b.cpu(),
		InFlight:  atomic.LoadInt64(&b.inflight),
		MaxFlight: 0,
		MinRt:     0,
		MaxPass:   0,
	}
}
