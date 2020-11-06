package cpu

import "sync/atomic"

var (
	usage uint64
)

type State struct {
	Usage uint64 // CPU 的使用率
}

func ReadState(state *State) {
	state.Usage = atomic.LoadUint64(&usage)
}
