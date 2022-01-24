package main

import(
	"fmt"
)

func main() {
	var q Queue

	q.Add(1)
	q.Add(2)
	q.Add(3)
	q.Add(4)
	for q.Next() {
		fmt.Println(q.Pop())
	}

	q.Add(5)
	q.Add(6)
	q.Add(7)
	q.Add(8)
	for q.Next() {
		fmt.Println(q.Pop())
	}
}

type Queue struct {
	begin *LinkedList
	end *LinkedList
}

func (q *Queue) Next() bool {
	if q.begin == nil {
		return false
	}
	return true
}

func (q *Queue) Add(value int) {
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
