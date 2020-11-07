package breaker

import "time"

type Breaker interface {
	Allow() error
	MakeSuccess() // 若 Allow() == nil 则 外界需要调用该函数
	MakeFailed()  // 若 Allow() != nil 则 外界需要调用该函数
}

type Config struct {
	SwitchOff bool    // 开关 true 为开启 默认为 false
	K         float64 // SRE 中的 盐值  可以根据实际情况 进行调整

	Window  time.Duration
	Bucket  int
	Request int64 // 触发SRE的最小 请求数
}

func (c *Config) fix() {
	if c.K == 0 {
		c.K = 1.5
	}

	if c.Request == 0 {
		c.Request = 100
	}

	if c.Bucket == 0 {
		c.Bucket = 10
	}

	if c.Window == 0 {
		c.Window = 3 * time.Second
	}
}

const (
	StateOpen     = iota // 表示开启
	StateClosed          // 表示关闭
	StateHalfOpen        // 表示 半开启
)

func newBreaker(c *Config) Breaker {
	return newSre(c)
}
