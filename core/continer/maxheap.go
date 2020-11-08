package continer

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	HEAP_CAP = 1 << 10
)

type Node struct {
	ele        unsafe.Pointer
	createTime int64
	idx        int64
}

func NewNode(ele unsafe.Pointer) *Node {
	return &Node{
		ele:        ele,
		createTime: time.Now().UnixNano(),
	}
}

func testNewNode(ele unsafe.Pointer, num int64) *Node {
	return &Node{
		ele:        ele,
		createTime: num,
	}
}

type MaxHeap struct {
	nodeList []*Node
	len      int64
	mu       sync.Mutex
	cap      int64
}

func NewMaxHeap(len int64) *MaxHeap {

	if len == 0 {
		len = HEAP_CAP
	}

	return &MaxHeap{
		nodeList: make([]*Node, len),
		len:      0,
		mu:       sync.Mutex{},
		cap:      len,
	}
}

func (h *MaxHeap) Push(node *Node) bool {
	if node == nil {
		return false
	}
	if h.isFull() {
		return false
	}
	pos := h.Len()

	h.mu.Lock()

	for ; pos > 1 && h.nodeList[pos/2] != nil && h.nodeList[pos/2].createTime < node.createTime; pos /= 2 {
		h.nodeList[pos] = h.nodeList[pos/2]
	}

	h.nodeList[pos] = node
	h.mu.Unlock()

	h.addLen()
	return true

}

func (h *MaxHeap) Pop() *Node {

	if h.isEmpty() {
		return nil
	}

	last := h.nodeList[h.Len()-1]
	h.mu.Lock()
	ele := h.nodeList[0]

	parent, child := 1, 0

	for ; int64(parent*2) < h.Len(); parent = child {
		child = parent * 2

		if int64(child) < h.Len() && h.nodeList[child].createTime < h.nodeList[child+1].createTime {
			child++
		}

		if last.createTime < h.nodeList[child].createTime {
			h.nodeList[parent] = h.nodeList[child]
		} else {
			break
		}

	}
	h.nodeList[parent] = last
	h.mu.Unlock()
	h.subLen()
	return ele
}

func (h *MaxHeap) Len() int64 {
	return atomic.LoadInt64(&h.len)
}

func (h *MaxHeap) isEmpty() bool {

	if atomic.CompareAndSwapInt64(&h.len, 0, h.len) {
		return true
	}
	return false
}

func (h *MaxHeap) isFull() bool {
	if atomic.CompareAndSwapInt64(&h.len, h.cap-1, h.len) {
		return true
	}

	return false
}

func (h *MaxHeap) addLen() {
	for atomic.CompareAndSwapInt64(&h.len, h.len, h.len+1) {
		return
	}
}

func (h *MaxHeap) subLen() {

	if h.Len() > 0 {
		for atomic.CompareAndSwapInt64(&h.len, h.len, h.len-1) {
			return
		}
	}

}
