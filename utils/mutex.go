package utils

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const mutextLocked = 1 << iota

type Mutexs struct {
	sync.Mutex
}

func (m *Mutexs) TryLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutextLocked)
}
