// TODO modify bit shift operator to pre-calculateted array for performance
package wm

import (
	"algo/heap"
	"algo/sds"
	"algo/que"
)

// WaveletMatrix is the struct of the Wavelet matrix.
// bitVectors is bits of the original slice.
// zeroNums is the number of zero of the bitsVector.
// firstIndexes is the first index of values in the final slice that is generated from bitVectors. 0-indexed.
type WaveletMatrix struct {
	bitVectors []*sds.SuccinctDictionary
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
	w := &WaveletMatrix{make([]*sds.SuccinctDictionary, length), make([]int, length), make(map[int]int)}

	s0 := make([]int, len(t)) // numbers of previous bit 0
	s1 := make([]int, 0) // numbers of previous bit 1
	copy(s0, t)

	setNext := func(n0, n1, s []int, bit, start int, sd *sds.SuccinctDictionary) ([]int, []int) {
		for i, v := range s {
			if v & (1<<bit) > 0 {
				n1 = append(n1, v)
				sd.Set(start + i, true)
			} else {
				n0 = append(n0, v)
			}
		}
		return n0, n1
	}

	for i:=topBit; i>=0; i-- {
		var n0, n1 []int // next numbers of previous bit 0 and 1
		sd := sds.NewSuccinctDictionary(len(t))
		n0, n1 = setNext(n0, n1, s0, i, 0, sd)
		n0, n1 = setNext(n0, n1, s1, i, len(s0), sd)
		s0 = n0
		s1 = n1
		sd.Build()
		w.bitVectors[i] = sd
		w.zeroNums[i] = sd.Rank0(sd.Size())
	}

	s := s0
	start := 0
	for i:=0; i<len(t); i++ {
		if i == len(s0) {
			s = s1
			start -= len(s0)
		}
		_, ok := w.firstIndexes[s[start + i]]
		if !ok {
			w.firstIndexes[s[start + i]] = i
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
