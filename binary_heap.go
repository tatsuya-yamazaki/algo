package main

import(
	"fmt"
)

type Heap struct {
	list []int
	isChild func(parent, child int) bool
}

const(
	ASCENDING = true
	DECENDING = false
)

func NewHeap(ascending bool) *Heap {
	h := &Heap{make([]int, 0), nil}
	if ascending {
		h.isChild = func(parent, child int) bool { return h.list[parent] < h.list[child] }
	} else {
		h.isChild = func(parent, child int) bool { return h.list[parent] > h.list[child] }
	}
	return h
}

func (*Heap) parent(index int) int {
	// index is zero base
	return (index - 1) / 2
}

func (*Heap) left(index int) int {
	// index is zero base
	return index * 2 + 1
}

func (*Heap) right(index int) int {
	// index is zero base
	return (index + 1) * 2
}

func (h *Heap) Push(value int) {
	h.list = append(h.list, value)
	index := len(h.list) - 1
	for index != 0 {
		parent := h.parent(index)
		if h.isChild(parent, index) {
			break
		}
		h.list[parent], h.list[index] = h.list[index], h.list[parent]
		index = parent
	}
}

func (h *Heap) Top() int {
	return h.list[0]
}

func (h *Heap) Pop() int {
	ret := h.list[0]
	h.list[0] = h.list[len(h.list)-1]
	h.list = h.list[:len(h.list)-1]
	index := 0
	for len(h.list) > h.left(index) {
		child, right := h.left(index), h.right(index)
		if len(h.list) > right && h.isChild(right, child) {
			child = right
		}
		if h.isChild(index, child) {
			break
		} else {
			h.list[index], h.list[child] = h.list[child], h.list[index]
			index = child
		}
	}
	return ret
}

func (h *Heap) Next() bool {
	return len(h.list) != 0
}

func main() {
	h := NewHeap(false)
	h.Push(3)
	fmt.Println(h)
	h.Push(1)
	fmt.Println(h)
	h.Push(2)
	fmt.Println(h)
	h.Pop()
	fmt.Println(h)
}
