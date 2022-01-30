package main

import(
	"fmt"
)

func main() {
	d := NewDeque()

	d.AddLeft(1)
	d.AddLeft(2)
	d.AddLeft(3)
	d.AddLeft(4)
	d.AddLeft(5)
	d.AddLeft(6)
	d.AddLeft(7)
	d.AddLeft(8)
	for d.Next() {
		fmt.Println(d.PopRight())
	}

	d.AddLeft(1)
	d.AddLeft(2)
	d.AddLeft(3)
	d.AddLeft(4)
	d.AddLeft(5)
	d.AddLeft(6)
	d.AddLeft(7)
	d.AddLeft(8)
	for d.Next() {
		fmt.Println(d.PopLeft())
	}

	d.AddRight(1)
	d.AddRight(2)
	d.AddRight(3)
	d.AddRight(4)
	d.AddRight(5)
	d.AddRight(6)
	d.AddRight(7)
	d.AddRight(8)
	for d.Next() {
		fmt.Println(d.PopRight())
	}

	d.AddRight(1)
	d.AddRight(2)
	d.AddRight(3)
	d.AddRight(4)
	d.AddRight(5)
	d.AddRight(6)
	d.AddRight(7)
	d.AddRight(8)
	for d.Next() {
		fmt.Println(d.PopLeft())
	}

	d.AddLeft(1)
	d.AddRight(2)
	d.AddLeft(3)
	d.AddRight(4)
	d.AddLeft(5)
	d.AddRight(6)
	d.AddLeft(7)
	d.AddRight(8)
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

func (d *Deque) AddLeft(value int) {
	ll := &LinkedList{nil, d.begin, value}
	if d.begin == nil {
		d.end = ll
	} else {
		d.begin.prev = ll
	}
	d.begin = ll
}

func (d *Deque) AddRight(value int) {
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
