package sds

// This array contains the sum of 1 bits from the beginning to each digit position in all uint8 values.
var bitCounts [256][8]int

func init() {
	// bitCounts initialization
	for i := 0; i < 256; i++ {
		var bitCount [8]int
		c := 0
		for j := 0; j < 8; j++ {
			if i&(1<<j) > 0 {
				c++
			}
			bitCount[j] = c
		}
		bitCounts[i] = bitCount
	}
}

type SuccinctDictionary struct {
	size  int
	large []int // max bits size N is 2^31 - 1 (max int32)
	small []uint16
	bits  []uint8
}

// SMALL_SIZE * m = LARGE_SIZE (m >= 2)
// BITS_SIZE * l = SMALL_SIZE (l >= 2)
const (
	LARGE_SIZE                  = 1024    // (log2(N+1))^2
	SMALL_SIZE                  = 16      // log2(N+1) / 2
	BITS_SIZE                   = 8       // uint8 size
	SELECT_BLOCK_SIZE           = 1024    // (log2(N+1))^2
	SELECT_BLOCK_TYPE_THRESHOLD = 1048576 // (log2(N+1))^4
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

		bitCount := bitCounts[v][BITS_SIZE-1]
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
		ret += int(bitCounts[s.bits[i]][BITS_SIZE-1])
	}

	ret += bitCounts[s.bits[bitsIndex]][index%BITS_SIZE]

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
