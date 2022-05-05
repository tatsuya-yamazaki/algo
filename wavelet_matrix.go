package main

import (
	"fmt"
)

func main() {
}

type WaveletMatrix struct {
	t []int
	[]*SuccinctDictionary
}

func NewWaveletMatrix(t []int) {
	if len(t) == 0 {
		return nil
	}
	max := 0
	for _, v := range t {
		if max < v {
			max = v
		}
	}
	nt := make([]int, len(t))
	copy(nt, t)
	w := &WaveletMatrix{nt, nil}

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

func (s SuccinctDictionary) Size() int {
	return s.size
}

// 0 origin
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
