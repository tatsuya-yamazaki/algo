package wm

import (
	"sort"
	"algo/heap"
	"algo/sds"
	"algo/que"
)

// WaveletMatrix is the struct of the Wavelet matrix.
// bitVectors is bits of the original slice.
// zeroNums is the number of zero of the bitsVector.
// firstIndexes is the first index of values in the final slice that is generated from bitVectors.
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

	type sortInt struct {
		v, b int // v is value, b is bit
	}
	sis := make([]sortInt, len(t))
	for i:=0; i<len(t); i++ {
		sis[i].v = t[i]
	}

	for i:=topBit; i>=0; i-- {
		b := sds.NewSuccinctDictionary(len(sis))
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

// Access returns original slice item value.
// index is 0-origin.
func (w WaveletMatrix) Access(index int) int {
	r := 0
	for i:=len(w.bitVectors)-1; i>=0; i-- {
		b := w.bitVectors[i]
		if b.Access(index) {
			r += 1<<i
			index = w.zeroNums[i] + b.Rank(index) - 1 // 0-origin
		} else {
			index = b.Rank0(index) - 1
		}
	}
	return r
}

// Rank returns number of value appeared in 0 to index in original slice.
// index is 0-origin.
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
			index = w.zeroNums[i] + rank - 1 // 0-origin
		} else {
			index = b.Rank0(index) - 1
			// No applicable data
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

// Select returns index of value appeared specified times from original slice index 0.
// rank is the ascending rank of the values in the array. 0-origin
func (w WaveletMatrix) Select(value, rank int) int {
	out := w.bitVectors[0].Size() //out of range
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

// Quantile returns nth smallest value in specified interval of the original array.
// l, r are half-open interval. ex) [0, 1)
// rank is the rank of values in the array in ascending order. 0-origin
func (w WaveletMatrix) Quantile(l, r, rank int) int {
	value := 0
	for i:=len(w.bitVectors)-1; i>=0; i-- {
		b := w.bitVectors[i]
		one := 0 // number of 1 in [l, r) of s
		rightOne := 0 // mumber of 1 in r) of s
		if r > 0 {
			rightOne = b.Rank(r - 1)
			one += rightOne
		}
		leftOne := 0 // mumber of 1 in l) of s
		if l > 0 {
			leftOne = b.Rank(l - 1)
			one -= leftOne
		}
		zero := r - l - one // number of 0 in [l, r) of s
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
// k is the number of items you want to be return. 1-origin.
func (w WaveletMatrix) Topk(l, r, k int) (ret [][2]int) {
	h := heap.NewHeap(heap.DESCENDING)
	bits := len(w.bitVectors)
	h.Add(topkNode{l, r, bits-1, 0})
	bv := make([]int, bits)
	for i:=0; i<bits; i++ {
		bv[i] = 1<<i
	}
	bl := w.bitVectors[0].Size() - 1 // bitVector last bit index
	for h.Next() && k > 0 {
		n := h.Pop().(topkNode)
		if n.i == -1 {
			k--
			ret = append(ret, [2]int{n.v, n.r - n.l})
			continue
		}
		b := w.bitVectors[n.i]
		one := 0 // num of 1 bit [l, r)
		if n.r > 0 {
			one += b.Rank(n.r-1)
		}
		leftOne := 0 // num of 1 bit l)
		leftZero := 0 // num of 0 bit l)
		if n.l > 0 {
			leftOne += b.Rank(n.l-1)
			one -= leftOne
			leftZero += b.Rank0(n.l-1)
		}
		zero := n.r - n.l - one // num of 0 bit [l, r)
		ni := n.i - 1 // new index of bitVector
		if zero > 0 {
			h.Add(topkNode{leftZero, leftZero + zero, ni, n.v})
		}
		if one > 0 {
			ol := b.Rank0(bl) + leftOne // new l of first 1 bit
			h.Add(topkNode{ol, ol+one, ni, n.v + bv[n.i]})
		}
	}
	return
}

// Topk returns sum of value in [l, r).
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
// l1, r1 are half-open interval. ex) [0, 1). 0-origin
// l2, r2 are half-open interval. ex) [0, 1). 0-origin
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
		// n1 are not 0, so v.r1 are not 0.
		one1 := b.Rank(v.r1 - 1) // number of one in v.r1)
		leftOne1 := 0 // number of one in v.l1)
		if v.l1 > 0 {
			leftOne1 += b.Rank(v.l1 - 1)
		}
		leftZero1 := v.l1 - leftOne1 // number of zero in v.l1)
		zero1 := v.r1 - one1 // number of zero in v.r1)

		// n2 are not 0, so v.r2 are not 0.
		one2 := b.Rank(v.r2 - 1) // number of one in v.r2)
		leftOne2 := 0 // number of one in v.l2)
		if v.l2 > 0 {
			leftOne2 += b.Rank(v.l2 - 1)
		}
		leftZero2 := v.l2 - leftOne2 // number of zero in v.l2)
		zero2 := v.r2 - one2 // number of zero in v.r2)

		zero := b.Rank0(b.Size()-1) // number of zero in b
		bit := 1 << v.i
		v.i-- // next index of bitVectors

		q.Add(intersectValue{leftZero1, zero1, leftZero2, zero2, v.i, v.v})
		q.Add(intersectValue{zero+leftOne1, zero+one1, zero+leftOne2, zero+one2, v.i, v.v+bit})
	}
	return
}
