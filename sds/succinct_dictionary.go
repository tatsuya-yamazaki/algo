package sds

type SuccinctDictionary struct {
	size  int
	large []int // max bits size N is 2**31 - 1 (max int32)
	small []uint16
	bits  []uint8
}

// SMALL_SIZE * m = LARGE_SIZE (m >= 2)
// BITS_SIZE * l = SMALL_SIZE (l >= 2)
const (
	LARGE_SIZE = 1024 // (log2(N+1))**2
	SMALL_SIZE = 16   // log2(N+1) / 2
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
	s.large = make([]int, getSuitableLength(LARGE_SIZE)+1)
	s.small = make([]uint16, getSuitableLength(SMALL_SIZE)+1)
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

func getLargeIndex(index int) int {
	return index / LARGE_SIZE
}

func getSmallIndex(index int) int {
	return index / SMALL_SIZE
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
	s.small[0] = 0
	s.large[0] = 0
	beforeSmallIndex := 0
	beforeLargeIndex := 0
	for i, v := range s.bits {
		index := i * BITS_SIZE
		smallIndex := getSmallIndex(index) + 1
		largeIndex := getLargeIndex(index) + 1

		if beforeSmallIndex < smallIndex {
			s.small[smallIndex] = s.small[beforeSmallIndex]
			beforeSmallIndex = smallIndex
		}

		if beforeLargeIndex < largeIndex {
			s.large[largeIndex] = s.large[beforeLargeIndex]
			s.small[smallIndex-1] = 0
			beforeLargeIndex = largeIndex
		}

		bitCount := bitNums[v]
		s.small[smallIndex] += uint16(bitCount)
		s.large[largeIndex] += int(bitCount)
	}
}

// Rank returns 1 bit num in [0, r)
func (s SuccinctDictionary) Rank(r int) (ret int) {
	if r < 1 {
		return 0
	}
	index := r - 1
	largeIndex := getLargeIndex(index)
	ret += int(s.large[largeIndex])

	smallIndex := getSmallIndex(index)
	ret += int(s.small[smallIndex])

	bitsIndex := getBitsIndex(index)
	for i := bitsIndex - 1; i >= 0 && smallIndex == getSmallIndex(i*BITS_SIZE); i-- {
		ret += int(bitNums[s.bits[i]])
	}

	bits := uint8(s.bits[bitsIndex])
	for i := uint8(1); i <= getBit(index) && i > 0; i <<= 1 {
		if i&bits > 0 {
			ret++
		}
	}

	return ret
}

// Rank returns 0 bit num in [0, r)
func (s SuccinctDictionary) Rank0(r int) int {
	if r < 1 {
		return 0
	}
	if r > s.size {
		return s.size
	}
	return r - s.Rank(r)
}

// Select returns index where 1 bit appears n times.
// The index is 0-indexed.
// If the number of 1 bits is less than n, it returns out of range index.
func (s SuccinctDictionary) Select(n int) int {
	l, r := 1, s.size+1
	for l < r {
		m := (l + r) / 2
		rank := s.Rank(m)
		if rank < n {
			l = m + 1
		} else {
			r = m
		}
	}
	return l - 1
}

// Select returns index where 0 bit appears n times.
// The index is 0-indexed.
// If the number of 0 bits is less than n, it returns out of range index.
func (s SuccinctDictionary) Select0(n int) int {
	l, r := 1, s.size+1
	for l < r {
		m := (l + r) / 2
		rank := s.Rank0(m)
		if rank < n {
			l = m + 1
		} else {
			r = m
		}
	}
	return l - 1
}
