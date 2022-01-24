package main

import(
	"fmt"
)

func main() {
	var q Queue

	q.add(1)
	q.add(2)
	q.add(3)
	q.add(4)
	q.add(5)
	q.add(6)
	q.add(7)
	q.add(8)

	for q.next() {
		fmt.Println(q.pop())
	}
}

type Queue struct {
	begin *LinkedList
	end *LinkedList
}

func (q *Queue) next() bool {
	if q.begin == nil {
		return false
	}
	return true
}

func (q *Queue) add(value int) {
	ll := &LinkedList{q.end, nil, value}
	if q.begin == nil {
		q.begin = ll
	}
	if q.end != nil {
		q.end.next = ll
	}
	q.end = ll
}

func (q *Queue) pop() int {
	value := q.begin.value
	if q.begin == q.end {
		q.end = nil
	}
	q.begin = q.begin.next
	return value
}

type LinkedList struct {
        prev, next *LinkedList
	value int
}
