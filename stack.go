package main

import(
	"fmt"
)

func main() {
	s := NewStack()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)
	s.Push(5)
	s.Push(6)
	s.Push(7)
	s.Push(8)
	fmt.Println(s.Top())
	for s.Next() {
		fmt.Println(s.Pop())
	}
}

type Stack struct {
	list []int
}

func NewStack() *Stack {
	return &Stack{make([]int, 0)}
}

func (s *Stack) Next() bool {
	return len(s.list) > 0
}

func (s *Stack) Top() int {
	return s.list[len(s.list)-1]
}

func (s *Stack) Push(x int) {
	s.list = append(s.list, x)
}

func (s *Stack) Pop() int {
	ret := s.Top()
	s.list = s.list[:len(s.list)-1]
	return ret
}
