// TODO If Atcoder adopts Go v1.8.0 or later, fix it to use generics.
// TODO Then, the wavelet matrix must be fixed.
package heap

// Node is the interface a node of Heap.
// Less returns Node is less than a or not.
// Greater returns Node is greater than a or not.
type Node interface {
	Less(a *Node) bool
	Greater(a *Node) bool
}

// Heap is the binary heap structure.
// To use it, the Node interface must be implemented.
// Its indexes are 0-origin.
// It can use ascending or descending order.
// TODO It may need to be refactored, expecially Pop().
// TODO It may need to be devided into min heap and max heap. Then remove isChild, use Less or Greater
type Heap struct {
	n []*Node
	isChild func(parent, child int) bool
}

const(
	ASCENDING = true
	DECENDING = false
)

func NewHeap(ascending bool) *Heap {
	h := &Heap{make([]*Node, 0), nil}
	if ascending {
		h.isChild = func(parent, child int) bool { return (*h.n[parent]).Less(h.n[child]) }
	} else {
		h.isChild = func(parent, child int) bool { return (*h.n[parent]).Greater(h.n[child]) }
	}
	return h
}

func parent(i int) int {
	return (i - 1) / 2
}

func left(i int) int {
	return i * 2 + 1
}

func right(i int) int {
	return (i + 1) * 2
}

func (h *Heap) Add(value Node) {
	h.n = append(h.n, &value)
	i := len(h.n) - 1
	for i != 0 {
		p := parent(i)
		if h.isChild(p, i) {
			break
		}
		h.n[p], h.n[i] = h.n[i], h.n[p]
		i =p
	}
}

func (h *Heap) Top() *Node {
	return h.n[0]
}

func (h *Heap) Pop() Node {
	ret := h.n[0]
	last := len(h.n)-1
	h.n[0] = h.n[last]
	h.n = h.n[:last]
	i := 0
	for last > left(i) {
		c, r := left(i), right(i)
		if len(h.n) > r && h.isChild(r, c) {
			c = r
		}
		if h.isChild(i, c) {
			break
		} else {
			h.n[i], h.n[c] = h.n[c], h.n[i]
			i = c
		}
	}
	return *ret
}

func (h *Heap) Next() bool {
	return len(h.n) != 0
}
