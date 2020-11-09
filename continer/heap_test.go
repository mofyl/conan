package continer

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestMaxHeap(t *testing.T) {
	h := NewMaxHeap(3)

	num1 := 1

	h.Push(testNewNode(unsafe.Pointer(&num1), 1))
	h.Push(testNewNode(unsafe.Pointer(&num1), 2))
	h.Push(testNewNode(unsafe.Pointer(&num1), 3))

	for i := 0; i < int(h.Len()); i++ {
		e := h.Pop()
		fmt.Println(e.createTime)
	}
}
