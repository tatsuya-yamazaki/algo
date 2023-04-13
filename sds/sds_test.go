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
			t.Errorf("%v", s.small)
			t.Errorf("%v", s.large)
			return
		}
	}
	// out of range
	for i := 1039; i <= 1040; i++ {
		if a := s.Rank(i); 1038 != a {
			t.Errorf("%v != %v", 1038, a)
			return
		}
	}
}

func TestRank0(t *testing.T) {
	s := NewSuccinctDictionary(1038)
	s.Build()
	for i := 0; i <= 1038; i++ {
		if a := s.Rank0(i); i != a {
			t.Errorf("%v != %v", i, a)
			return
		}
	}
	// out of range
	for i := 1039; i <= 1040; i++ {
		if a := s.Rank0(i); 1038 != a {
			t.Errorf("%v != %v", 1038, a)
			return
		}
	}
}

func TestSelect(t *testing.T) {
	s := NewSuccinctDictionary(1038)
	for i := 0; i < 1038; i++ {
		s.Set(i, true)
	}
	s.Build()
	if a := s.Select(0); 0 != a {
		t.Errorf("%v != %v", 0, a)
		return
	}
	for i := 1; i <= 1038; i++ {
		if a := s.Select(i); i-1 != a {
			t.Errorf("%v != %v", i-1, a)
			return
		}
	}
	// out of range
	for i := 1039; i <= 1040; i++ {
		if a := s.Select(i); 1038 != a {
			t.Errorf("%v != %v", 1038, a)
			return
		}
	}
}

func TestSelect0(t *testing.T) {
	s := NewSuccinctDictionary(1038)
	s.Build()
	if a := s.Select0(0); 0 != a {
		t.Errorf("%v != %v", 0, a)
		return
	}
	for i := 1; i <= 1038; i++ {
		if a := s.Select0(i); i-1 != a {
			t.Errorf("%v != %v", i-1, a)
			return
		}
	}
	// out of range
	for i := 1039; i <= 1040; i++ {
		if a := s.Select0(i); 1038 != a {
			t.Errorf("%v != %v", 1038, a)
			return
		}
	}
}

func BenchmarkNewSuccinctDictionary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewSuccinctDictionary(200000)
	}
}

func BenchmarkBuild(b *testing.B) {
	s := NewSuccinctDictionary(200000)
	for i := 0; i < 200000; i++ {
		if rand.Intn(2) > 0 {
			s.Set(i, true)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Build()
	}
}

func BenchmarkSet(b *testing.B) {
	s := NewSuccinctDictionary(200000)
	for i := 0; i < 200000; i++ {
		if rand.Intn(2) > 0 {
			s.Set(i, true)
		}
	}
	s.Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := rand.Intn(200000)
		b := false
		if rand.Intn(2) > 0 {
			b = true
		}
		s.Set(j, b)
	}
}

func BenchmarkAccess(b *testing.B) {
	s := NewSuccinctDictionary(200000)
	for i := 0; i < 200000; i++ {
		if rand.Intn(2) > 0 {
			s.Set(i, true)
		}
	}
	s.Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Access(rand.Intn(200000))
	}
}

func BenchmarkRank(b *testing.B) {
	s := NewSuccinctDictionary(200000)
	for i := 0; i < 200000; i++ {
		if rand.Intn(2) > 0 {
			s.Set(i, true)
		}
	}
	s.Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Rank(rand.Intn(200000))
	}
}

func BenchmarkRank0(b *testing.B) {
	s := NewSuccinctDictionary(200000)
	for i := 0; i < 200000; i++ {
		if rand.Intn(2) > 0 {
			s.Set(i, true)
		}
	}
	s.Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Rank0(rand.Intn(200000))
	}
}

func BenchmarkSelect(b *testing.B) {
	s := NewSuccinctDictionary(200000)
	c := 0
	for i := 0; i < 200000; i++ {
		if rand.Intn(2) > 0 {
			c++
			s.Set(i, true)
		}
	}
	s.Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Select(rand.Intn(c + 1))
	}
}

func BenchmarkSelect0(b *testing.B) {
	s := NewSuccinctDictionary(200000)
	c := 200000
	for i := 0; i < 200000; i++ {
		if rand.Intn(2) > 0 {
			c--
			s.Set(i, true)
		}
	}
	s.Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Select0(rand.Intn(c + 1))
	}
}
