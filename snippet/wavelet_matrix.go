
// Node is the interface a node of Heap.
// Less returns Node is less than a or not.
// Greater returns Node is greater than a or not.
type HeapNode interface {
	Less(a *HeapNode) bool
	Greater(a *HeapNode) bool
}

// Heap is the binary heap structure.
// To use it, the Node interface must be implemented.
// Its indexes are 0-origin.
// It can use ascending or descending order.
// TODO It may need to be refactored, expecially Pop().
// TODO It may need to be devided into min heap and max  Then remove isChild, use Less or Greater
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

type SuccinctDictionary struct {
	size int
	chunks []int // max bits size N is 2**31 - 1 (max int32)
	blocks []uint16
	bits   []uint8
}

// BLOCK_SIZE * m = CHUNK_SIZE (m >= 2)
// BITS_SIZE * l = BLOCK_SIZE (l >= 2)
const (
	CHUNK_SIZE = 1024 // (log2(N+1))**2
	BLOCK_SIZE = 16   // log2(N+1) / 2
	BITS_SIZE  = 8    // uint8 size
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

// index is 0-indexed
func (s SuccinctDictionary) Access(index int) bool {
	b := s.bits[getBitsIndex(index)]
	return b&getBit(index) > 0
}

// index is 0-indexed
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

// Rank returns 1 bit num in [0, r)
func (s SuccinctDictionary) Rank(r int) (ret int) {
	if r < 1 {
		return 0
	}
	index := r - 1
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

// Rank returns 0 bit num in [0, r)
func (s SuccinctDictionary) Rank0(r int) int {
	if r < 1 {
		return 0
	}
	return r - s.Rank(r)
}

// Select returns index where 1 bit appears n times.
// The index is 1-indexed.
// If the real index is i, it returns i + 1.
func (s SuccinctDictionary) Select(n int) int {
	l, r := 0, s.size
	for l < r {
		m := (l + r) / 2
		rank := s.Rank(m)
		if rank < n {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

// Select returns index where 0 bit appears n times.
// The index is 1-indexed.
// If the real index is i, it returns i + 1.
func (s SuccinctDictionary) Select0(n int) int {
	l, r := 0, s.size
	for l < r {
		m := (l + r) / 2
		rank := s.Rank0(m)
		if rank < n {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

// WaveletMatrix is the struct of the Wavelet matrix.
// bitVectors is bits of the original slice.
// zeroNums is the number of zero of the bitsVector.
// firstIndexes is the first index of values in the final slice that is generated from bitVectors. 0-indexed.
type WaveletMatrix struct {
	bitVectors []*SuccinctDictionary
	zeroNums []int
	firstIndexes map[int]int
}

// NewWaveletMatrix returns pointer of WaveletMatrix.
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
		v, b int // v is value, b is bit
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
		w.zeroNums[i] = b.Rank0(b.Size())
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

// Access returns original slice item value.
// index is 0-indexed.
func (w WaveletMatrix) Access(index int) int {
	index++ // fix to 1-indexed
	value := 0
	for i:=len(w.bitVectors)-1; i>=0; i-- {
		b := w.bitVectors[i]
		if b.Access(index - 1) {
			value += 1<<i
			index = w.zeroNums[i] + b.Rank(index)
		} else {
			index = b.Rank0(index)
		}
	}
	return value
}

// Rank returns number of value appeared the interval [0, index) in original slice.
func (w WaveletMatrix) Rank(value, index int) int {
	fi, ok := w.firstIndexes[value]
	if !ok {
		return 0
	}
	for i:=len(w.bitVectors)-1; i>=0; i-- {
		b := w.bitVectors[i]
		if value & (1<<i) > 0 {
			rank := b.Rank(index)
			// No applicable data
			if rank == 0 {
				return 0
			}
			index = w.zeroNums[i] + rank // 1-indexed
		} else {
			index = b.Rank0(index) // 1-indexed
			// No applicable data
			if index == 0 {
				return 0
			}
		}
	}
	ret := index - fi
	if index < 0 {
		return 0
	} else {
		return ret
	}
}

// Select returns index of value appeared specified times from original slice. 1-indexed.
// rank is the ascending rank of the value in the array. 1-indexed.
func (w WaveletMatrix) Select(value, rank int) int {
	last := w.bitVectors[0].Size()
	fi, ok := w.firstIndexes[value]
	index := fi + rank
	if !ok || rank < 1 || last < index || w.Rank(value, last) < rank {
		return 0
	}

	for i:=0; i<len(w.bitVectors); i++ {
		b := w.bitVectors[i]
		if value & (1<<i) > 0 {
			index = b.Select(index - w.zeroNums[i])
		} else {
			index = b.Select0(index)
		}
	}
	return index
}

// Quantile returns nth smallest value in specified interval of the original array.
// l, r are half-open interval. ex) [0, 1)
// rank is the rank of values in the array in ascending order. 1-indexed
func (w WaveletMatrix) Quantile(l, r, rank int) int {
	value := 0
	for i:=len(w.bitVectors)-1; i>=0; i-- {
		b := w.bitVectors[i]
		rightOne := b.Rank(r) // number of 1 in r) of s
		leftOne := b.Rank(l) // number of 1 in l) of s
		one := rightOne - leftOne // number of 1 in [l, r) of s
		zero := r - l - one // number of 0 in [l, r) of s
		if rank > zero {
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

// topkNode is used by Topk priority queue
// It implements HeapNode.
// l, r are half-open interval. ex) [0, 1)
// i is the index of bitVectors.
// v is the accumulative value of bit.
type topkNode struct {
	l, r, i, v int
}

// Less returns whether n is less than a or not.
func (n topkNode) Less(a *HeapNode) bool {
	v := (*a).(topkNode)
	return (n.r - n.l) < (v.r - v.l)
}

// Less returns whether n is greater than a or not.
func (n topkNode) Greater(a *HeapNode) bool {
	v := (*a).(topkNode)
	return (n.r - n.l) > (v.r - v.l)
}

// Topk returns top k frequent values in [l, r).
// return array is sort by frequency in descending order,
// but is not stable original order.
// l, r are half-open interval. ex) [0, 1).
// k is the number of items you want to be return. 1-indexed.
func (w WaveletMatrix) Topk(l, r, k int) (ret [][2]int) {
	h := NewHeap(DESCENDING)
	bits := len(w.bitVectors)
	h.Add(topkNode{l, r, bits-1, 0})
	bv := make([]int, bits)
	for i:=0; i<bits; i++ {
		bv[i] = 1<<i
	}
	for h.Next() && k > 0 {
		n := h.Pop().(topkNode)
		if n.i == -1 {
			k--
			ret = append(ret, [2]int{n.v, n.r - n.l})
			continue
		}
		b := w.bitVectors[n.i]
		leftOne := b.Rank(n.l) // num of 1 bit l)
		leftZero := n.l - leftOne // num of 0 bit l)
		one := b.Rank(n.r) - leftOne // num of 1 bit [l, r)
		zero := n.r - n.l - one // num of 0 bit [l, r)
		ni := n.i - 1 // new index of bitVector
		if zero > 0 {
			h.Add(topkNode{leftZero, leftZero + zero, ni, n.v})
		}
		if one > 0 {
			ol := w.zeroNums[n.i] + leftOne // new l of first 1 bit
			h.Add(topkNode{ol, ol+one, ni, n.v + bv[n.i]})
		}
	}
	return
}

// Sum returns sum of value in [l, r).
// l, r are half-open interval. ex) [0, 1).
func (w WaveletMatrix) Sum(l, r int) (ret int) {
	k := r - l
	for _, v := range w.Topk(l, r, k) {
		ret += v[0] * v[1]
	}
	return
}

// intersectNode is used by Intersect queue
// It implements QueueValue.
// l1, r1 are half-open interval. ex) [0, 1).
// l2, r2 are half-open interval. ex) [0, 1).
// i is the index of bitVectors.
// v is the accumulative value of bit.
type intersectValue struct {
	l1, r1, l2, r2, i, v int
}

// Intersect returns the common values and their frequency in [l1, r1) and [l2, r2).
// l1, r1 are half-open interval. ex) [0, 1). 0-indexed
// l2, r2 are half-open interval. ex) [0, 1). 0-indexed
func (w WaveletMatrix) Intersect(l1, r1, l2, r2 int) (ret [][3]int) {
	q := NewQueue()
	q.Add(intersectValue{l1, r1, l2, r2, len(w.bitVectors)-1, 0})
	for q.Next() {
		v := q.Pop().(intersectValue)
		n1 := v.r1 - v.l1 // length of [l1, r1)
		n2 := v.r2 - v.l2 // length of [l2, r2)
		// If there are no common values.
		if n1 == 0 || n2 == 0 {
			continue
		}
		if v.i == -1 {
			ret = append(ret, [3]int{v.v, n1, n2})
			continue
		}

		b := w.bitVectors[v.i]
		one1 := b.Rank(v.r1) // number of one in v.r1)
		leftOne1 := b.Rank(v.l1) // number of one in v.l1)
		leftZero1 := v.l1 - leftOne1 // number of zero in v.l1)
		zero1 := v.r1 - one1 // number of zero in v.r1)

		one2 := b.Rank(v.r2) // number of one in v.r2)
		leftOne2 := b.Rank(v.l2) // number of one in v.l2)
		leftZero2 := v.l2 - leftOne2 // number of zero in v.l2)
		zero2 := v.r2 - one2 // number of zero in v.r2)

		zero := w.zeroNums[v.i] // number of zero in b
		bit := 1 << v.i
		v.i-- // next index of bitVectors

		q.Add(intersectValue{leftZero1, zero1, leftZero2, zero2, v.i, v.v})
		q.Add(intersectValue{zero+leftOne1, zero+one1, zero+leftOne2, zero+one2, v.i, v.v+bit})
	}
	return
}

// Rangefreq returns the number of value between x and y - 1 in the interval [l, r) of the original array.
// l, r are half-open interval. ex) [0, 1).
// The values is greater than and equal x.
// The values is less than y.
func (w WaveletMatrix) Rangefreq(l, r, x, y int) (ret int) {
	if l >= r || x >= y {
		return 0
	}
	y--
	var s [][3]int // node of slice. index 0: node.l, 1: node.r, 2: node.value
	s = append(s, [3]int{l, r, 0})
	var xv, yv int
	for i:=0; i<63; i++ {
		if i < len(w.bitVectors) {
			continue
		}
		bit := 1<<i
		xv += x & bit
		yv += y & bit
	}
	for i:=len(w.bitVectors)-1; i>=0; i-- {
		bit := 1<<i
		xv += x & bit
		yv += y & bit
		bv := w.bitVectors[i]
		z := w.zeroNums[i]
		isRange := func(v int) bool {
			if v < xv || v > yv {
				return false
			}
			return true
		}
		var ns [][3]int // next node slice
		add := func(s [][3]int, l, r, v int) [][3]int {
			n := r - l
			if n > 0 {
				if i == 0 {
					ret += n
				} else {
					s = append(s, [3]int{l, r, v})
				}
			}
			return s
		}
		for _, v := range s {
			zv := v[2]
			ov := v[2] + bit
			rol := bv.Rank(v[0]) // number of 1 in node.l)
			rzl := v[0] - rol // number of 0 in node.l)
			ror := bv.Rank(v[1]) // number of 1 in node.r)
			rzr := v[1] - ror // number of 0 in node.r)
			if isRange(zv) {
				ns = add(ns, rzl, rzr, zv)
			}
			if isRange(ov) {
				ns = add(ns, z + rol, z + ror, ov)
			}
		}
		s = ns
	}
	return ret
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
	ll := &queueLinkedList{nil, value}
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
		q.begin = q.begin.next
	}
	return value
}

type QueueValue interface {
}

type queueLinkedList struct {
        next *queueLinkedList
	value QueueValue
}
