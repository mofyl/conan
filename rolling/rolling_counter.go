package rolling

import "time"

type RollingCounter interface {
	Metric
	Aggregation
	Timespan() int
	Reduce(f func(iterator Iterator) float64) float64
}

type RollingCounterOpts struct {
	Size           int
	BucketDuration time.Duration
}

type rollingCounter struct {
	policy *Policy
}

func NewRollingCounter(opts RollingCounterOpts) RollingCounter {
	window := NewWindow(WindowOpt{Size: opts.Size})
	p := NewPolicy(window, PolicyOpts{BucketDuration: opts.BucketDuration})
	return &rollingCounter{
		policy: p,
	}
}

func (p *rollingCounter) Reduce(f func(iterator Iterator) float64) float64 {
	return p.policy.Reduce(f)
}

func (p *rollingCounter) Add(val float64) {
	if val < 0 {
		return
	}
	p.policy.Add(val)
}

func (p *rollingCounter) Timespan() int {
	return p.policy.timespan()
}

/*
	Min() float64
	Max() float64
	Avg() float64
	Sum() float64
*/
func (p *rollingCounter) Avg() float64 {
	return p.policy.Reduce(Avg)
}

func (p *rollingCounter) Min() float64 {
	return p.policy.Reduce(Min)
}

func (p *rollingCounter) Max() float64 {
	return p.policy.Reduce(Max)
}

func (p *rollingCounter) Sum() float64 {
	return p.policy.Reduce(Sum)
}

func (p *rollingCounter) Value() float64 {
	return p.Sum()
}
