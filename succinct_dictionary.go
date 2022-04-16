package main
import "fmt"

func main() {
	s := NewSuccinctDictionary(17)
	s.b[0] = 128
	fmt.Println(s.Access(7))
}

type SuccinctDictionary struct {
	l []int
	s []uint8
	b []uint64
}

func NewSuccinctDictionary(n int) *SuccinctDictionary {
	ret := &SuccinctDictionary{}
	num := n / 64
	if n % 64 > 0 {
		num++
	}
	ret.l = make([]int, num / 2 + num % 2)
	ret.s = make([]uint8, num)
	ret.b = make([]uint64, num)
	return ret
}

func (SuccinctDictionary) getBit(i int) uint64 {
	return 1<<(i % 64)
}

func (s SuccinctDictionary) Access(i int) bool {
	b := s.b[i / 64]
	return b & s.getBit(i) > 0
}

func (s *SuccinctDictionary) Set(i int, b bool) {
	if s.Access(i) == b {
		return
	}
	index := i / 64
	bit := s.getBit(i)
	if b {
		s.b[index] = s.b[index] + bit
		s.s[index]++
	} else {
		s.b[index] = s.b[index] - bit
		s.s[index]--
	}
}
