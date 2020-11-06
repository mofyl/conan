package bbr

import "context"

type Op int


type DoneInfo struct {
	Err error
	Op Op
}


type Limiter interface {
	Allow(ctx context.Context)(func (DoneInfo)  ,error)
}


