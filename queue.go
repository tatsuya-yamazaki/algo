package main

import(
	"fmt"
)

func main() {
	q := NewQueue()

	q.Push(1)
	q.Push(2)
	q.Push(3)
	q.Push(4)
	for q.Next() {
		fmt.Println(q.Pop())
	}

	q.Push(5)
	q.Push(6)
	q.Push(7)
	q.Push(8)
	for q.Next() {
		fmt.Println(q.Pop())
	}
}

type Queue struct {
	begin *LinkedList
	end *LinkedList
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

func (q *Queue) Push(value int) {
	ll := &LinkedList{q.end, nil, value}
	if q.end == nil {
		q.begin = ll
	} else {
		q.end.next = ll
	}
	q.end = ll
}

func (q *Queue) Pop() int {
	value := q.begin.value
	if q.begin == q.end {
		q.begin = nil
		q.end = nil
	} else {
		q.begin.next.prev = nil
		q.begin = q.begin.next
	}
	return value
}

type LinkedList struct {
        prev, next *LinkedList
	value int
}
