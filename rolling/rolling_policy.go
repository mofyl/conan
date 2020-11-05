package rolling

import (
	"fmt"
	"sync"
	"time"
)

type Policy struct {
	mu     *sync.Mutex
	size   int
	window *Window
	offset int

	bucketDuration time.Duration
	lastAppendTime time.Time
}

type PolicyOpts struct {
	BucketDuration time.Duration
}

func NewPolicy(window *Window, opts PolicyOpts) *Policy {
	return &Policy{
		mu:             &sync.Mutex{},
		size:           window.Size(),
		window:         window,
		bucketDuration: opts.BucketDuration,
		lastAppendTime: time.Now(),
		offset:         0,
	}
}

// 计算从上次添加 到现在一共 经历了多少个 bucket
func (p *Policy) timespan() int {
	t := time.Since(p.lastAppendTime)
	fmt.Println(t)
	return int(t / p.bucketDuration)
}

// 主要用于计算 当前的offset。 并且清空 这段时间经过的bucket
func (p *Policy) add(f func(offset int, val float64), val float64) {
	p.mu.Lock()
	timespan := p.timespan()
	fmt.Println(timespan)
	if timespan > 0 {
		// 更新appendTime
		p.lastAppendTime = p.lastAppendTime.Add(time.Duration(timespan * int(p.bucketDuration)))
		// 计算当前应该插入的位置
		s := p.offset + 1
		if timespan > p.size {
			timespan = p.size
		}
		e, e1 := s+timespan, 0
		if e > p.size {
			e1 = e - p.size
			e = p.size
		}
		offset := p.offset

		for i := s; i < e; i++ {
			// 这里就是 将 这段时间内经过的 bucket 清空
			p.window.ResetBucket(i)
			offset = i
		}
		// 注意这里是从 0开始的
		// 这里从0开始是有原因的
		for i := 0; i < e1; i++ {
			// 这里从头开始清空
			p.window.ResetBucket(i)
			offset = i
		}
		p.offset = offset
	}
	f(p.offset, val)
	p.mu.Unlock()
}

func (p *Policy) Append(val float64) {
	p.add(p.window.Append, val)
}

func (p *Policy) Add(val float64) {
	p.add(p.window.Add, val)
}

// 将f 应用到  now - lastAppendTime 之间经过的bucket之外的 bucket的point 上
func (p *Policy) Reduce(f func(iterator Iterator) float64) (val float64) {
	p.mu.Lock()
	timespan := p.timespan()
	if count := p.size - timespan; count > 0 {
		offset := p.offset + 1 + timespan
		if offset >= p.size {
			offset = offset - p.size
		}
		val = f(p.window.Iterator(offset, count))
	}
	p.mu.Unlock()
	return
}
