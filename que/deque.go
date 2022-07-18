package que

type Deque struct {
	begin *dequeLinkedList
	end   *dequeLinkedList
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

func (d *Deque) AddLeft(value DequeValue) {
	ll := &dequeLinkedList{nil, d.begin, value}
	if d.begin == nil {
		d.end = ll
	} else {
		d.begin.prev = ll
	}
	d.begin = ll
}

func (d *Deque) AddRight(value DequeValue) {
	ll := &dequeLinkedList{d.end, nil, value}
	if d.end == nil {
		d.begin = ll
	} else {
		d.end.next = ll
	}
	d.end = ll
}

func (d *Deque) PopLeft() DequeValue {
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

func (d *Deque) PopRight() DequeValue {
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

type DequeValue interface {
}

type dequeLinkedList struct {
	prev, next *dequeLinkedList
	value      DequeValue
}
