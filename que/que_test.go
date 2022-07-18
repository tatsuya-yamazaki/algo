package que

import (
	"reflect"
	"testing"
)

type queueValue struct {
	v int
}

type dequeValue struct {
	v int
}

func TestQueue(t *testing.T) {
	e := []queueValue{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}}
	q := NewQueue()
	for _, v := range e {
		q.Add(v)
	}

	a := make([]queueValue, 0)
	for q.Next() {
		a = append(a, q.Pop().(queueValue))
	}

	if !reflect.DeepEqual(e, a) {
		t.Errorf("%v != %v", e, a)
	}
}

func TestDequeAddLeftPopRight(t *testing.T) {
	e := []dequeValue{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}}
	d := NewDeque()
	for _, v := range e {
		d.AddLeft(v)
	}

	a := make([]dequeValue, 0)
	for d.Next() {
		a = append(a, d.PopRight().(dequeValue))
	}

	if !reflect.DeepEqual(e, a) {
		t.Errorf("%v != %v", e, a)
	}
}

func TestDequeAddLeftPopLeft(t *testing.T) {
	e := []dequeValue{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}}
	d := NewDeque()
	for _, v := range e {
		d.AddLeft(v)
	}

	a := make([]dequeValue, len(e))
	i := len(e)
	for d.Next() {
		i--
		a[i] = d.PopLeft().(dequeValue)
	}

	if !reflect.DeepEqual(e, a) {
		t.Errorf("%v != %v", e, a)
	}
}

func TestDequeAddRightPopLeft(t *testing.T) {
	e := []dequeValue{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}}
	d := NewDeque()
	for _, v := range e {
		d.AddRight(v)
	}

	a := make([]dequeValue, 0)
	for d.Next() {
		a = append(a, d.PopLeft().(dequeValue))
	}

	if !reflect.DeepEqual(e, a) {
		t.Errorf("%v != %v", e, a)
	}
}

func TestDequeAddRightPopRight(t *testing.T) {
	e := []dequeValue{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}}
	d := NewDeque()
	for _, v := range e {
		d.AddRight(v)
	}

	a := make([]dequeValue, len(e))
	i := len(e)
	for d.Next() {
		i--
		a[i] = d.PopRight().(dequeValue)
	}

	if !reflect.DeepEqual(e, a) {
		t.Errorf("%v != %v", e, a)
	}
}
