package breaker

import (
	"errors"
	"lb/rolling"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// SRE 阈值 requests ==  k * accepts  requests 小于该值的时候 不用熔断。 大于开始熔断
// 熔断概率p =  max(0 , (requests - k * accepts) / requests+1)
type sre struct {
	stat   rolling.RollingCounter
	randMu sync.Mutex // 这里 random.Rand 不是线程安全的 所以要给锁
	r      *rand.Rand

	k       float64
	request int64
	state   int32
}

func newSre(conf *Config) *sre {
	couterOpt := rolling.RollingCounterOpts{
		Size:           conf.Bucket,
		BucketDuration: time.Duration(int64(conf.Window) / int64(conf.Bucket)),
	}
	stat := rolling.NewRollingCounter(couterOpt)

	return &sre{
		stat:    stat,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
		k:       conf.K,
		request: conf.Request,
		state:   StateClosed,
	}
}

func (s *sre) summary() (success int64, total int64) {
	s.stat.Reduce(func(iterator rolling.Iterator) float64 {
		for iterator.Next() {
			b := iterator.Bucket()
			total += b.Count

			for _, p := range b.Point {
				success += int64(p)
			}
		}
		return 0
	})
	return
}

// 随机值 在 [0,p) 之间 表示需要丢弃  大于p表示不需要丢弃
// true 表示需要丢弃 false 表示不需要丢弃
func (s *sre) turnOnSre(p float64) (turnOn bool) {
	s.randMu.Lock()
	turnOn = s.r.Float64() < p
	s.randMu.Unlock()
	return
}

func (s *sre) Allow() error {
	success, total := s.summary()

	k := float64(success) * s.k

	// 这里表示不需要熔断 或者结束熔断
	if total < s.request || float64(total) < k {
		if atomic.LoadInt32(&s.state) == StateOpen {
			atomic.CompareAndSwapInt32(&s.state, StateOpen, StateClosed)
		}
		return nil
	}

	// 这里表示需要打开
	if atomic.LoadInt32(&s.state) == StateClosed {
		atomic.CompareAndSwapInt32(&s.state, StateClosed, StateOpen)
	}

	p := math.Max(0, (float64(total)-k)/float64(total+1))

	if s.turnOnSre(p) {
		// 返回 503 过载保护
		return errors.New("service unavailable")
	}
	return nil
}

func (s *sre) MakeSuccess() {
	s.stat.Add(1)
}

func (s *sre) MakeFailed() {
	// 这里表示 本次需要熔断 上报的就是0 。在summary中 加success的值的时候 就会加到0
	// 但是Count 的值还是会累计
	s.stat.Add(0)
}
