package bbr

import "context"

const (
	Success Op = iota
	Ignore
	Drop
)

type Op int

type DoneInfo struct {
	Err error
	Op  Op
}

type Limiter interface {
	Allow(ctx context.Context) (func(DoneInfo), error)
}
