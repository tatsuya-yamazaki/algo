package que

import(
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
	e := []queueValue{{1},{2},{3},{4},{5},{6},{7},{8}}
	q := NewQueue()
	for _, v := range e {
		q.Push(v)
	}

	a := make([]queueValue, 0)
	for q.Next() {
		a = append(a, q.Pop().(queueValue))
	}

	if !reflect.DeepEqual(e, a) {
		t.Errorf("%v != %v", e, a)
	}
}

func TestDequePushLeftPopRight(t *testing.T) {
	e := []dequeValue{{1},{2},{3},{4},{5},{6},{7},{8}}
	d := NewDeque()
	for _, v := range e {
		d.PushLeft(v)
	}

	a := make([]dequeValue, 0)
	for d.Next() {
		a = append(a, d.PopRight().(dequeValue))
	}

	if !reflect.DeepEqual(e, a) {
		t.Errorf("%v != %v", e, a)
	}
}

func TestDequePushLeftPopLeft(t *testing.T) {
	e := []dequeValue{{1},{2},{3},{4},{5},{6},{7},{8}}
	d := NewDeque()
	for _, v := range e {
		d.PushLeft(v)
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

func TestDequePushRightPopLeft(t *testing.T) {
	e := []dequeValue{{1},{2},{3},{4},{5},{6},{7},{8}}
	d := NewDeque()
	for _, v := range e {
		d.PushRight(v)
	}

	a := make([]dequeValue, 0)
	for d.Next() {
		a = append(a, d.PopLeft().(dequeValue))
	}

	if !reflect.DeepEqual(e, a) {
		t.Errorf("%v != %v", e, a)
	}
}

func TestDequePushRightPopRight(t *testing.T) {
	e := []dequeValue{{1},{2},{3},{4},{5},{6},{7},{8}}
	d := NewDeque()
	for _, v := range e {
		d.PushRight(v)
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
