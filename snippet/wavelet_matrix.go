
type SuccinctDictionary struct {
	size int
	chunks []int
	blocks []uint16
	bits   []uint8
}

const (
	CHUNK_SIZE = 1024
	BLOCK_SIZE = 16
	BITS_SIZE  = 8
)

func NewSuccinctDictionary(size int) *SuccinctDictionary {
	if size <= 0 || size >= (1<<31) {
		return nil
	}
	s := &SuccinctDictionary{}
	s.size = size
	getSuitableLength := func(n int) int {
		ret := size / n
		if size%n > 0 {
			ret++
		}
		return ret
	}
	s.chunks = make([]int, getSuitableLength(CHUNK_SIZE))
	s.blocks = make([]uint16, getSuitableLength(BLOCK_SIZE))
	s.bits = make([]uint8, getSuitableLength(BITS_SIZE))
	return s
}

var bitNums = [256]uint8{
	0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	4, 5, 5, 6, 5, 6, 6, 7, 5, 6, 6, 7, 6, 7, 7, 8,
}

func getBit(n int) uint8 {
	return 1 << (n % BITS_SIZE)
}

func getChunkIndex(index int) int {
	return index / CHUNK_SIZE
}

func getBlockIndex(index int) int {
	return index / BLOCK_SIZE
}

func getBitsIndex(index int) int {
	return index / BITS_SIZE
}

func (s SuccinctDictionary) Size() int {
	return s.size
}

func (s SuccinctDictionary) Access(index int) bool {
	b := s.bits[getBitsIndex(index)]
	return b&getBit(index) > 0
}

func (s *SuccinctDictionary) Set(index int, b bool) {
	if b == s.Access(index) {
		return
	}
	bit := getBit(index)
	bits := &s.bits[getBitsIndex(index)]
	if b {
		*bits += bit
		return
	}
	*bits -= bit
}

func (s *SuccinctDictionary) Build() {
	s.chunks[0] = 0
	s.blocks[0] = 0
	ci, bi := 0, 0
	for i, v := range s.bits {
		index := i * BITS_SIZE
		cin := getChunkIndex(index)
		bin := getBlockIndex(index)
		if ci < cin {
			s.chunks[cin] = s.chunks[ci]
			ci = cin
			s.blocks[bin] = 0
			bi = bin
		}
		if bi < bin {
			s.blocks[bin] = s.blocks[bi]
			bi = bin
		}
		c := bitNums[v]
		s.chunks[ci] += int(c)
		s.blocks[bi] += uint16(c)
	}
}

func (s SuccinctDictionary) Rank(index int) (ret int) {
	chunkIndex := getChunkIndex(index)
	if chunkIndex > 0 {
		ret += int(s.chunks[chunkIndex-1])
	}

	blockIndex := getBlockIndex(index)
	if blockIndex > 0 && (BLOCK_SIZE * blockIndex % CHUNK_SIZE != 0) {
		ret += int(s.blocks[blockIndex-1])
	}

	bitsIndex := getBitsIndex(index)
	bits := uint8(s.bits[bitsIndex])
	for i := uint8(1); (i <= getBit(index) && i > 0); i <<= 1 {
		if i&bits > 0 {
			ret++
		}
	}

	for i := bitsIndex - 1; i >= 0 && blockIndex == getBlockIndex(i*BITS_SIZE); i-- {
		ret += int(bitNums[s.bits[i]])
	}

	return ret
}

func (s SuccinctDictionary) Select(n int) int {
	l, r := 0, s.size
	var m int
	for l < r {
		m = (l + r) / 2
		rank := s.Rank(m)
		if rank < n {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

func (s SuccinctDictionary) Rank0(index int) int {
	return index + 1 - s.Rank(index)
}

func (s SuccinctDictionary) Select0(n int) int {
	l, r := 0, s.size
	var m int
	for l < r {
		m = (l + r) / 2
		rank := s.Rank0(m)
		if rank < n {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

type Queue struct {
	begin *queueLinkedList
	end *queueLinkedList
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Next() bool {
	if q.begin == nil {
		return false
	}
	return true
}

func (q *Queue) Add(value QueueValue) {
	ll := &queueLinkedList{q.end, nil, value}
	if q.end == nil {
		q.begin = ll
	} else {
		q.end.next = ll
	}
	q.end = ll
}

func (q *Queue) Pop() QueueValue {
	value := q.begin.value
	if q.begin == q.end {
		q.begin = nil
		q.end = nil
	} else {
		q.begin.next.prev = nil
		q.begin = q.begin.next
	}
	return value
}

type QueueValue interface {
}

type queueLinkedList struct {
        prev, next *queueLinkedList
	value QueueValue
}

type HeapNode interface {
	Less(a *HeapNode) bool
	Greater(a *HeapNode) bool
}

type Heap struct {
	n []*HeapNode
	isChild func(parent, child int) bool
}

const(
	ASCENDING = true
	DESCENDING = false
)

func NewHeap(ascending bool) *Heap {
	h := &Heap{make([]*HeapNode, 0), nil}
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

func (h *Heap) Add(value HeapNode) {
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

func (h *Heap) Top() *HeapNode {
	return h.n[0]
}

func (h *Heap) Pop() HeapNode {
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

type WaveletMatrix struct {
	bitVectors []*SuccinctDictionary
	zeroNums []int
	firstIndexes map[int]int
}

func NewWaveletMatrix(t []int) *WaveletMatrix {
	if len(t) == 0 {
		return nil
	}
	max := 0
	for i, v := range t {
		if t[max] < v {
			max = i
		}
	}
	topBit := 0
	for i:=0; i<64; i++ {
		if t[max] & (1<<i) > 0 {
			topBit = i
		}
	}

	length := topBit + 1
	w := &WaveletMatrix{make([]*SuccinctDictionary, length), make([]int, length), make(map[int]int)}

	type sortInt struct {
		v, b int
	}
	sis := make([]sortInt, len(t))
	for i:=0; i<len(t); i++ {
		sis[i].v = t[i]
	}

	for i:=topBit; i>=0; i-- {
		b := NewSuccinctDictionary(len(sis))
		for j, v := range sis {
			if v.v & (1<<i) > 0 {
				b.Set(j, true)
				sis[j].b = 1
			} else {
				sis[j].b = 0
			}
		}
		b.Build()
		w.bitVectors[i] = b
		w.zeroNums[i] = b.Size() - b.Rank(b.Size()-1)
		sort.SliceStable(sis, func(k, l int) bool {return sis[k].b < sis[l].b})
	}
	for i:=0; i<len(sis); i++ {
		_, ok := w.firstIndexes[sis[i].v]
		if !ok {
			w.firstIndexes[sis[i].v] = i
		}
	}
	return w
}

func (w WaveletMatrix) Access(index int) int {
	r := 0
	for i:=len(w.bitVectors)-1; i>=0; i-- {
		b := w.bitVectors[i]
		if b.Access(index) {
			r += 1<<i
			index = w.zeroNums[i] + b.Rank(index) - 1
		} else {
			index = b.Rank0(index) - 1
		}
	}
	return r
}

func (w WaveletMatrix) Rank(value, index int) int {
	fi, ok := w.firstIndexes[value]
	if !ok {
		return 0
	}
	for i:=len(w.bitVectors)-1; i>=0; i-- {
		b := w.bitVectors[i]
		if value & (1<<i) > 0 {
			rank := b.Rank(index)
			if rank == 0 {
				return 0
			}
			index = w.zeroNums[i] + rank - 1
		} else {
			index = b.Rank0(index) - 1
			if index < 0 {
				return 0
			}
		}
	}
	if index < fi {
		return 0
	} else {
		return index - fi + 1
	}
}

func (w WaveletMatrix) Select(value, rank int) int {
	out := w.bitVectors[0].Size()
	fi, ok := w.firstIndexes[value]
	index := fi + rank
	if !ok || rank < 0 || out <= index {
		return out
	}

	for i:=0; i<len(w.bitVectors); i++ {
		b := w.bitVectors[i]
		if value & (1<<i) > 0 {
			index = b.Select(index + 1 - w.zeroNums[i])
		} else {
			index = b.Select0(index + 1)
		}
		if out <= index {
			return out
		}
	}
	if value == w.Access(index) {
		return index
	}
	return out
}

func (w WaveletMatrix) Quantile(l, r, rank int) int {
	value := 0
	for i:=len(w.bitVectors)-1; i>=0; i-- {
		b := w.bitVectors[i]
		one := 0
		rightOne := 0
		if r > 0 {
			rightOne = b.Rank(r - 1)
			one += rightOne
		}
		leftOne := 0
		if l > 0 {
			leftOne = b.Rank(l - 1)
			one -= leftOne
		}
		zero := r - l - one
		if rank + 1 > zero {
			value += 1<<i
			z := w.zeroNums[i]
			l = z + leftOne
			r = z + rightOne
			rank = rank - zero
		} else {
			l = l - leftOne
			r = r - rightOne
		}
	}
	return value
}

type topkNode struct {
	l, r, i, v int
}

func (n topkNode) Less(a *HeapNode) bool {
	v := (*a).(topkNode)
	return (n.r - n.l) < (v.r - v.l)
}

func (n topkNode) Greater(a *HeapNode) bool {
	v := (*a).(topkNode)
	return (n.r - n.l) > (v.r - v.l)
}

func (w WaveletMatrix) Topk(l, r, k int) (ret [][2]int) {
	h := NewHeap(DESCENDING)
	bits := len(w.bitVectors)
	h.Add(topkNode{l, r, bits-1, 0})
	bv := make([]int, bits)
	for i:=0; i<bits; i++ {
		bv[i] = 1<<i
	}
	bl := w.bitVectors[0].Size() - 1
	for h.Next() && k > 0 {
		n := h.Pop().(topkNode)
		if n.i == -1 {
			k--
			ret = append(ret, [2]int{n.v, n.r - n.l})
			continue
		}
		b := w.bitVectors[n.i]
		one := 0
		if n.r > 0 {
			one += b.Rank(n.r-1)
		}
		leftOne := 0
		leftZero := 0
		if n.l > 0 {
			leftOne += b.Rank(n.l-1)
			one -= leftOne
			leftZero += b.Rank0(n.l-1)
		}
		zero := n.r - n.l - one
		ni := n.i - 1
		if zero > 0 {
			h.Add(topkNode{leftZero, leftZero + zero, ni, n.v})
		}
		if one > 0 {
			ol := b.Rank0(bl) + leftOne
			h.Add(topkNode{ol, ol+one, ni, n.v + bv[n.i]})
		}
	}
	return
}

func (w WaveletMatrix) Sum(l, r int) (ret int) {
	k := r - l
	for _, v := range w.Topk(l, r, k) {
		ret += v[0] * v[1]
	}
	return
}

type intersectValue struct {
	l1, r1, l2, r2, i, v int
}

func (w WaveletMatrix) Intersect(l1, r1, l2, r2 int) (ret [][3]int) {
	q := NewQueue()
	q.Add(intersectValue{l1, r1, l2, r2, len(w.bitVectors)-1, 0})
	for q.Next() {
		v := q.Pop().(intersectValue)
		n1 := v.r1 - v.l1
		n2 := v.r2 - v.l2
		if n1 == 0 || n2 == 0 {
			continue
		}
		if v.i == -1 {
			ret = append(ret, [3]int{v.v, n1, n2})
			continue
		}

		b := w.bitVectors[v.i]
		one1 := b.Rank(v.r1 - 1)
		leftOne1 := 0
		if v.l1 > 0 {
			leftOne1 += b.Rank(v.l1 - 1)
		}
		leftZero1 := v.l1 - leftOne1
		zero1 := v.r1 - one1

		one2 := b.Rank(v.r2 - 1)
		leftOne2 := 0
		if v.l2 > 0 {
			leftOne2 += b.Rank(v.l2 - 1)
		}
		leftZero2 := v.l2 - leftOne2
		zero2 := v.r2 - one2

		zero := b.Rank0(b.Size()-1)
		bit := 1 << v.i
		v.i--

		q.Add(intersectValue{leftZero1, zero1, leftZero2, zero2, v.i, v.v})
		q.Add(intersectValue{zero+leftOne1, zero+one1, zero+leftOne2, zero+one2, v.i, v.v+bit})
	}
	return
}
