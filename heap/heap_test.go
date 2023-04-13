package heap

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

type node struct {
	a int
}

func (n node) Less(a *HeapNode) bool {
	v := (*a).(node)
	return n.a < v.a
}

func TestAcsending(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var te []int
	ha := NewHeap(ASCENDING)
	for i := 0; i < 100; i++ {
		v := rand.Intn(100)
		te = append(te, v)
		ha.Add(node{v})
	}

	sort.Ints(te)
	var ta []int
	for _, v := range te {
		a := ha.Pop().(node).a
		ta = append(ta, a)
		if v != a {
			t.Errorf("%v != %v", v, a)
			t.Errorf("%v != %v", te, ta)
			break
		}
	}
}

func TestDescsending(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var te []int
	hd := NewHeap(DESCENDING)
	for i := 0; i < 100; i++ {
		v := rand.Intn(100)
		te = append(te, v)
		hd.Add(node{v})
	}

	sort.Sort(sort.Reverse(sort.IntSlice(te)))
	var ta []int
	for _, v := range te {
		a := hd.Pop().(node).a
		ta = append(ta, a)
		if v != a {
			t.Errorf("%v != %v", v, a)
			t.Errorf("%v != %v", te, ta)
			break
		}
	}
}
