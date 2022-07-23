// TODO randomize test data
// TODO add TestSet
package sds

import (
	"math/rand"
	"testing"
)

func TestSetAccess(t *testing.T) {
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
	for i := 0; i < 1038; i++ {
		s.Set(i, true)
	}
	s.Build()
	for i := 0; i <= 1038; i++ {
		if a := s.Rank(i); i != a {
			t.Errorf("%v != %v", i, a)
		}
	}
}

func TestRank0(t *testing.T) {
	s := NewSuccinctDictionary(1038)
	s.Build()
	for i := 0; i <= 1038; i++ {
		if a := s.Rank0(i); i != a {
			t.Errorf("%v != %v", i, a)
		}
	}
}

func TestSelect(t *testing.T) {
	s := NewSuccinctDictionary(1038)
	for i := 0; i < 1038; i++ {
		s.Set(i, true)
	}
	s.Build()
	for i := 0; i <= 1038; i++ {
		if a := s.Select(i); i != a {
			t.Errorf("%v != %v", i, a)
		}
	}
}

func TestSelect0(t *testing.T) {
	s := NewSuccinctDictionary(1038)
	s.Build()
	for i := 0; i <= 1038; i++ {
		if a := s.Select0(i); i != a {
			t.Errorf("%v != %v", i, a)
		}
	}
}

func BenchmarkSelect(b *testing.B) {
	s := NewSuccinctDictionary(200000)
	c := 0
	for i:=0; i<200000; i++ {
		if rand.Intn(1) > 0 {
			c++
			s.Set(i, true)
		}
	}
	s.Build()
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		s.Select(rand.Intn(c+1))
	}
}

func BenchmarkSelect0(b *testing.B) {
	s := NewSuccinctDictionary(200000)
	c := 200000
	for i:=0; i<200000; i++ {
		if rand.Intn(1) > 0 {
			c--
			s.Set(i, true)
		}
	}
	s.Build()
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		s.Select0(rand.Intn(c+1))
	}
}
