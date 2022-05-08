//TODO add TestSet
package sds

import (
	"testing"
)

func TestAccess(t *testing.T) {
	s := NewSuccinctDictionary(8)
	s.Set(7, true)
	s.Build()
	if e, a := true, s.Access(7); e != a {
		t.Errorf("%v != %v", e, a)
	}

	s.Set(7, true)
	s.Build()
	if e, a := true, s.Access(7); e != a {
		t.Errorf("%v != %v", e, a)
	}

	s.Set(7, false)
	s.Build()
	if e, a := false, s.Access(7); e != a {
		t.Errorf("%v != %v", e, a)
	}

	s.Set(7, false)
	s.Build()
	if e, a := false, s.Access(7); e != a {
		t.Errorf("%v != %v", e, a)
	}
}

func TestRank(t *testing.T) {
	s := NewSuccinctDictionary(1038)
	for i:=0; i<1038; i++ {
		s.Set(i, true)
	}
	s.Build()
	for i:=0; i<1038; i++ {
		if a := s.Rank(i); i + 1 != a {
			t.Errorf("%v != %v", i + 1, a)
		}
	}
}

func TestRank0(t *testing.T) {
	s := NewSuccinctDictionary(1038)
	s.Build()
	for i:=0; i<1038; i++ {
		if a := s.Rank0(i); i + 1 != a {
			t.Errorf("%v != %v", i + 1, a)
		}
	}
}

func TestSelect(t *testing.T) {
	s := NewSuccinctDictionary(1038)
	for i:=0; i<1038; i++ {
		s.Set(i, true)
	}
	s.Build()
	for i:=0; i<1038; i++ {
		if a := s.Select(i + 1); i != a {
			t.Errorf("%v != %v", i, a)
		}
	}
}

func TestSelect0(t *testing.T) {
	s := NewSuccinctDictionary(1038)
	s.Build()
	for i:=0; i<1038; i++ {
		if a := s.Select0(i + 1); i != a {
			t.Errorf("%v != %v", i, a)
		}
	}
}
