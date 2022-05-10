package main

import(
	"fmt"
)

func main() {
	d := NewDeque()

	d.PushLeft(1)
	d.PushLeft(2)
	d.PushLeft(3)
	d.PushLeft(4)
	d.PushLeft(5)
	d.PushLeft(6)
	d.PushLeft(7)
	d.PushLeft(8)
	for d.Next() {
		fmt.Println(d.PopRight())
	}

	d.PushLeft(1)
	d.PushLeft(2)
	d.PushLeft(3)
	d.PushLeft(4)
	d.PushLeft(5)
	d.PushLeft(6)
	d.PushLeft(7)
	d.PushLeft(8)
	for d.Next() {
		fmt.Println(d.PopLeft())
	}

	d.PushRight(1)
	d.PushRight(2)
	d.PushRight(3)
	d.PushRight(4)
	d.PushRight(5)
	d.PushRight(6)
	d.PushRight(7)
	d.PushRight(8)
	for d.Next() {
		fmt.Println(d.PopRight())
	}

	d.PushRight(1)
	d.PushRight(2)
	d.PushRight(3)
	d.PushRight(4)
	d.PushRight(5)
	d.PushRight(6)
	d.PushRight(7)
	d.PushRight(8)
	for d.Next() {
		fmt.Println(d.PopLeft())
	}

	d.PushLeft(1)
	d.PushRight(2)
	d.PushLeft(3)
	d.PushRight(4)
	d.PushLeft(5)
	d.PushRight(6)
	d.PushLeft(7)
	d.PushRight(8)
	for d.Next() {
		fmt.Println(d.PopLeft())
	}

}

type Deque struct {
	begin *LinkedList
	end *LinkedList
}

func NewDeque() *Deque {
	return &Deque{}
}

func (d *Deque) Next() bool {
	if d.begin == nil {
		return false
	}
	return true
}

func (d *Deque) PushLeft(value int) {
	ll := &LinkedList{nil, d.begin, value}
	if d.begin == nil {
		d.end = ll
	} else {
		d.begin.prev = ll
	}
	d.begin = ll
}

func (d *Deque) PushRight(value int) {
	ll := &LinkedList{d.end, nil, value}
	if d.end == nil {
		d.begin = ll
	} else {
		d.end.next = ll
	}
	d.end = ll
}

func (d *Deque) PopLeft() int {
	value := d.begin.value
	if d.begin == d.end {
		d.begin = nil
		d.end = nil
	} else {
		d.begin.next.prev = nil
		d.begin = d.begin.next
	}
	return value
}

func (d *Deque) PopRight() int {
	value := d.end.value
	if d.begin == d.end {
		d.begin = nil
		d.end = nil
	} else {
		d.end.prev.next = nil
		d.end = d.end.prev
	}
	return value
}

type LinkedList struct {
        prev, next *LinkedList
	value int
}
