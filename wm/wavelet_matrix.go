package wm

import (
	"algo/heap"
	"algo/que"
	"algo/sds"
)

var bits = [63]int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152, 4194304, 8388608, 16777216, 33554432, 67108864, 134217728, 268435456, 536870912, 1073741824, 2147483648, 4294967296, 8589934592, 17179869184, 34359738368, 68719476736, 137438953472, 274877906944, 549755813888, 1099511627776, 2199023255552, 4398046511104, 8796093022208, 17592186044416, 35184372088832, 70368744177664, 140737488355328, 281474976710656, 562949953421312, 1125899906842624, 2251799813685248, 4503599627370496, 9007199254740992, 18014398509481984, 36028797018963968, 72057594037927936, 144115188075855872, 288230376151711744, 576460752303423488, 1152921504606846976, 2305843009213693952, 4611686018427387904}

// WaveletMatrix is the struct of the Wavelet matrix.
// bitVectors is bits of the original slice.
// zeroNums is the number of zero of the bitsVector.
// firstIndexes is the first index of values in the final slice that is generated from bitVectors. 0-indexed.
type WaveletMatrix struct {
	bitVectors   []*sds.SuccinctDictionary
	zeroNums     []int
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
	for i := 0; i < 63; i++ {
		if t[max]&(bits[i]) > 0 {
			topBit = i
		}
	}

	length := topBit + 1
	w := &WaveletMatrix{make([]*sds.SuccinctDictionary, length), make([]int, length), make(map[int]int)}

	s := make([]int, len(t)) // previous numbers
	copy(s, t)
	ns := make([]int, len(t)) // next numbers
	p0, p1 := len(t), 0       // previous length of 0, 1
	n0, n1 := len(t), 0       // next length of 0, 1

	setNext := func(i, j, start int, sd *sds.SuccinctDictionary) {
		if s[j]&bits[i] > 0 {
			ns[len(t)-1-n1] = s[j] // set 1 bit number in reverse order
			sd.Set(n0+n1, true)
			n1++
		} else {
			ns[n0] = s[j]
			n0++
		}
	}

	for i := topBit; i >= 0; i-- {
		n0, n1 = 0, 0
		sd := sds.NewSuccinctDictionary(len(t))
		for j := 0; j < p0; j++ {
			setNext(i, j, 0, sd)
		}
		for j := len(t) - 1; j >= len(t)-p1; j-- {
			setNext(i, j, p0, sd)
		}
		sd.Build()
		w.bitVectors[i] = sd
		w.zeroNums[i] = sd.Rank0(sd.Size())
		s, ns = ns, s
		p0, p1 = n0, n1
	}

	var prev int // previous value
	setFirstIndex := func(i, start, j int) {
		v := s[i]
		if prev == v {
			return
		}
		prev = v
		w.firstIndexes[v] = start + j
	}

	if p0 > 0 {
		prev = s[0] + 1 // set a value different from the first term
		for i := 0; i < p0; i++ {
			setFirstIndex(i, 0, i)
		}
	}
	if p1 > 0 {
		prev = s[len(t)-1] + 1 // set a value different from the first term
		// 1 bit number is reverse order
		for i, j := len(t)-1, 0; i >= len(t)-p1; i, j = i-1, j+1 {
			setFirstIndex(i, p0, j)
		}
	}
	return w
}

// Top returns top bit index in original slice values.
// The return value is 0-origin.
func (w WaveletMatrix) Top() int {
	return len(w.bitVectors) - 1
}

// Access returns original slice item value.
// index is 0-indexed.
func (w WaveletMatrix) Access(index int) int {
	index++ // fix to 1-indexed
	value := 0
	for i := w.Top(); i >= 0; i-- {
		b := w.bitVectors[i]
		if b.Access(index - 1) {
			value += bits[i]
			index = w.zeroNums[i] + b.Rank(index)
		} else {
			index = b.Rank0(index)
		}
	}
	return value
}

// Rank returns number of values appeared the interval [0, index) in original slice.
func (w WaveletMatrix) Rank(value, index int) int {
	fi, ok := w.firstIndexes[value]
	if !ok {
		return 0
	}
	for i := w.Top(); i >= 0; i-- {
		b := w.bitVectors[i]
		if value&(bits[i]) > 0 {
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

// RankLess returns number of values are less than value in the interval [l, r) of the original slice.
func (w WaveletMatrix) RankLess(l, r, value int) (ret int) {
	if top := w.Top(); top < 62 && value >= bits[top]*2 {
		return r - l
	}
	for i := w.Top(); i >= 0; i-- {
		b := w.bitVectors[i]
		if value&(bits[i]) > 0 {
			rankLeft := b.Rank(l)
			one := b.Rank(r) - rankLeft
			ret += r - l - one
			l = w.zeroNums[i] + rankLeft
			r = l + one
		} else {
			l = b.Rank0(l)
			r = b.Rank0(r)
		}
	}
	return ret
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

	for i := 0; i <= w.Top(); i++ {
		b := w.bitVectors[i]
		if value&(bits[i]) > 0 {
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
	for i := w.Top(); i >= 0; i-- {
		b := w.bitVectors[i]
		rightOne := b.Rank(r)     // number of 1 in r) of s
		leftOne := b.Rank(l)      // number of 1 in l) of s
		one := rightOne - leftOne // number of 1 in [l, r) of s
		zero := r - l - one       // number of 0 in [l, r) of s
		if rank > zero {
			value += bits[i]
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
// It implements heap.HeapNode.
// l, r are half-open interval. ex) [0, 1)
// i is the index of bitVectors.
// v is the accumulative value of bit.
type topkNode struct {
	l, r, i, v int
}

// Less returns whether n is less than a or not.
func (n topkNode) Less(a *heap.HeapNode) bool {
	v := (*a).(topkNode)
	return (n.r - n.l) < (v.r - v.l)
}

// Less returns whether n is greater than a or not.
func (n topkNode) Greater(a *heap.HeapNode) bool {
	v := (*a).(topkNode)
	return (n.r - n.l) > (v.r - v.l)
}

// Topk returns top k frequent values in [l, r).
// return array is sort by frequency in descending order,
// but is not stable original order.
// l, r are half-open interval. ex) [0, 1).
// k is the number of items you want to be return. 1-indexed.
func (w WaveletMatrix) Topk(l, r, k int) (ret [][2]int) {
	h := heap.NewHeap(heap.DESCENDING)
	h.Add(topkNode{l, r, w.Top(), 0})
	for h.Next() && k > 0 {
		n := h.Pop().(topkNode)
		if n.i == -1 {
			k--
			ret = append(ret, [2]int{n.v, n.r - n.l})
			continue
		}
		b := w.bitVectors[n.i]
		leftOne := b.Rank(n.l)       // num of 1 bit l)
		leftZero := n.l - leftOne    // num of 0 bit l)
		one := b.Rank(n.r) - leftOne // num of 1 bit [l, r)
		zero := n.r - n.l - one      // num of 0 bit [l, r)
		ni := n.i - 1                // new index of bitVector
		if zero > 0 {
			h.Add(topkNode{leftZero, leftZero + zero, ni, n.v})
		}
		if one > 0 {
			ol := w.zeroNums[n.i] + leftOne // new l of first 1 bit
			h.Add(topkNode{ol, ol + one, ni, n.v + bits[n.i]})
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
// It implements que.QueueValue.
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
	q := que.NewQueue()
	q.Add(intersectValue{l1, r1, l2, r2, w.Top(), 0})
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
		one1 := b.Rank(v.r1)         // number of one in v.r1)
		leftOne1 := b.Rank(v.l1)     // number of one in v.l1)
		leftZero1 := v.l1 - leftOne1 // number of zero in v.l1)
		zero1 := v.r1 - one1         // number of zero in v.r1)

		one2 := b.Rank(v.r2)         // number of one in v.r2)
		leftOne2 := b.Rank(v.l2)     // number of one in v.l2)
		leftZero2 := v.l2 - leftOne2 // number of zero in v.l2)
		zero2 := v.r2 - one2         // number of zero in v.r2)

		zero := w.zeroNums[v.i] // number of zero in b
		bit := bits[v.i]
		v.i-- // next index of bitVectors

		q.Add(intersectValue{leftZero1, zero1, leftZero2, zero2, v.i, v.v})
		q.Add(intersectValue{zero + leftOne1, zero + one1, zero + leftOne2, zero + one2, v.i, v.v + bit})
	}
	return
}

// Rangefreq returns the number of value between x and y - 1 in the interval [l, r) of the original array.
// l, r are half-open interval. ex) [0, 1).
// The values is greater than and equal x.
// The values is less than y.
func (w WaveletMatrix) Rangefreq(l, r, x, y int) (ret int) {
	return w.RankLess(l, r, y) - w.RankLess(l, r, x)
}
