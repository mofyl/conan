package leaky

import (
	"sync/atomic"
	"time"
)

var (
	defaultConf = &LeakyConfig{
		GenerateTimeImMs: 1000,
	}
)

const (
	OPEN_LEAKY  = 1
	CLOSE_LEAKY = 2
)

type LeakyConfig struct {
	GenerateTimeImMs int32
	Cap              int64
}

type Leaky struct {
	tokenChan chan struct{}
	conf      *LeakyConfig
	closed    chan struct{}
	isClosed  uint32 // 1 表示正常运行  2 表示已经关闭
	count     uint64
	enable    uint32
}

func NewLeaky(config *LeakyConfig) *Leaky {

	if config == nil {
		config = defaultConf
	}
	l := &Leaky{
		tokenChan: make(chan struct{}),
		closed:    make(chan struct{}),
		conf:      config,
	}

	go l.generateToken()
	return l
}

func (l *Leaky) generateToken() {

	atomic.StoreUint32(&l.isClosed, 1)
	ticker := time.NewTicker(time.Duration(l.conf.GenerateTimeImMs) * time.Millisecond)
	for {
		select {
		case _, ok := <-l.closed:
			if !ok {
				return
			}
			ticker.Stop()
			atomic.StoreUint32(&l.isClosed, 2)
			close(l.tokenChan)

			return
		case <-ticker.C:
			if l.canWriteToken() {
				l.tokenChan <- struct{}{}
			}
		default:

		}
	}

}

func (l *Leaky) canWriteToken() bool {

	if atomic.LoadUint32(&l.isClosed) == 2 {
		return false
	}
	if atomic.LoadUint32(&l.enable) == CLOSE_LEAKY {
		return false
	}

	return true
}

func (l *Leaky) GenerateTime() int32 {
	return l.conf.GenerateTimeImMs
}

func (l *Leaky) SetEnable(enable uint32) bool {

	if atomic.LoadUint32(&l.isClosed) == 2 {
		return false
	}
	atomic.StoreUint32(&l.enable, enable)
	return true
}

// 返回true表示获取到了Token
// 返回false表示没有获取到Token
func (l *Leaky) TryGetToken() bool {
	_, ok := <-l.tokenChan
	if !ok {
		return false
	}
	return true
}

func (l *Leaky) Close() {
	close(l.closed)
}
