package main

import (
	"fmt"
	"sort"
)

func main() {
	t := []int{5,4,5,5,2,1,5,6,1,3,5,0}
	w := NewWaveletMatrix(t)
	m := make(map[int]struct{})
	for _, v := range t {
		m[v] = struct{}{}
	}

	fmt.Println("/////////////bit vectors")
	for i:=2; i>=0; i-- {
		for j:=0; j<len(t); j++ {
			if w.bitVectors[i].Access(j) {
				fmt.Print(1)
			} else {
				fmt.Print(0)
			}
		}
		fmt.Println()
	}

	fmt.Println("/////////////access/////////////")
	fmt.Println(t)
	for i:=0; i<len(t); i++ {
		fmt.Print(w.Access(i))
	}
	fmt.Println()

	fmt.Println("/////////////rank/////////////")
	for k, _ := range m {
		fmt.Println("value", k)
		fmt.Println(t)
		for i:=0; i<len(t); i++ {
			fmt.Print(w.Rank(k, i))
			fmt.Print(" ")
		}
		fmt.Println()
	}

	fmt.Println("/////////////select/////////////")
	for k, _ := range m {
		fmt.Println("value", k)
		fmt.Println(t)
		for i:=0; i<len(t); i++ {
			fmt.Print(w.Select(k, i))
			fmt.Print(" ")
		}
		fmt.Println()
	}

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
		s := NewSuccinctDictionary(len(sis))
		for j, v := range sis {
			if v.v & (1<<i) > 0 {
				s.Set(j, true)
				sis[j].b = 1
			} else {
				sis[j].b = 0
			}
		}
		s.Build()
		w.bitVectors[i] = s
		w.zeroNums[i] = s.Size() - s.Rank(s.Size()-1)
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
		s := w.bitVectors[i]
		if s.Access(index) {
			r += 1<<i
			index = w.zeroNums[i] + s.Rank(index) - 1 // 0-origin
		} else {
			index = s.Rank0(index) - 1
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
		s := w.bitVectors[i]
		if value & (1<<i) > 0 {
			rank := s.Rank(index)
			// No applicable data
			if rank == 0 {
				return 0
			}
			index = w.zeroNums[i] + rank - 1 // 0-origin
		} else {
			index = s.Rank0(index) - 1
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

// rank is the ascending rank of the values in the array. 0-origin
func (w WaveletMatrix) Select(value, rank int) int {
	out := w.bitVectors[0].Size() //out of range
	fi, ok := w.firstIndexes[value]
	index := fi + rank
	if !ok || rank < 0 || out <= index {
		return out
	}

	for i:=0; i<len(w.bitVectors); i++ {
		s := w.bitVectors[i]
		if value & (1<<i) > 0 {
			index = s.Select(index + 1 - w.zeroNums[i])
		} else {
			index = s.Select0(index + 1)
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

// return value is 0-origin index
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

// index is 0-origin
func (s SuccinctDictionary) Rank0(index int) int {
	return index + 1 - s.Rank(index)
}

// return value is 0-origin index
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
