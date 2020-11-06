package bbr

import (
	"lb/rolling"
	cpustate "lb/sys/cpu"
	"sync/atomic"
	"time"
)

var (
	cpu int64
	decay =0.95 // 这里这个值是看自己项目进行调整的
	initTime int64
	defaultConf = &Config{
		Window: time.Second * 10 ,
		WinBucket: 100,
		CPUThreshold: 800,
	}
)

func Init() {
	go cpuProc()
}

type cpuGetter func() int64


func cpuProc () {
	ticker := time.Ticker{}
	for range ticker.C{
		s := cpustate.State{}
		cpustate.ReadState(&s)
		prevCpu := atomic.LoadInt64(&cpu)
		// EWMA
		curCpu := int64( float64(prevCpu) * decay + (1.0-decay) * float64(s.Usage) )
		atomic.StoreInt64(&cpu , curCpu)
	}
}

type Config struct {
	Enable bool
	Window time.Duration
	WinBucket int
	CPUThreshold int64
}

type BBR struct {
	cpu cpuGetter
	passStat rolling.RollingCounter // 该窗口用于统计通过的 请求数量。也就是qps
	rtStat rolling.RollingCounter // 该窗口用于统计 每次请求的rtt
	inflight int64 // 当前正在请求的数量
	winBucketPreSec int64
	config *Config
	preDrop atomic.Value // 上次丢弃的时间 。这里的时间指 time.Since(上次计算时间) 这里保存的是一个time.Duration
}
