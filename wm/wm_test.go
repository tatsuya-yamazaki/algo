// TODO randomize test data. create function generateRandomData.
package wm

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestBits(t *testing.T) {
	for i := 0; i < 63; i++ {
		e := 1 << i
		a := bits[i]
		if e != a {
			t.Errorf("i == %v", i)
			t.Errorf("%v != %v", e, a)
		}
	}
}

func TestTop(t *testing.T) {
	s := []int{5, 4, 5, 5, 2, 1, 5, 6, 1, 3, 5, 0}
	e := -1
	for _, v := range s {
		for i := 0; i < 63; i++ {
			if v&(1<<i) > 0 {
				if i > e {
					e = i
				}
			}
		}
	}
	// If s is all zero, top bit is 1.
	if e == -1 {
		e = 0
	}
	w := NewWaveletMatrix(s)
	if a := w.Top(); e != a {
		t.Errorf("%v != %v", e, a)
	}
}

func TestAccess(t *testing.T) {
	s := []int{5, 4, 5, 5, 2, 1, 5, 6, 1, 3, 5, 0}
	w := NewWaveletMatrix(s)
	for i := 0; i < len(s); i++ {
		if a := w.Access(i); s[i] != a {
			t.Errorf("%v != %v", s[i], a)
		}
	}
}

func TestAllZero(t *testing.T) {
	s := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	w := NewWaveletMatrix(s)
	for i := 0; i < len(s); i++ {
		if a := w.Access(i); s[i] != a {
			t.Errorf("%v != %v", s[i], a)
		}
	}
	if e, a := 1, len(w.bitVectors); e != a {
		t.Errorf("bitVector number %v != %v", e, a)
	}
	if e, a := 0, w.Top(); e != a {
		t.Errorf("top bit %v != %v", e, a)
	}
}

func TestRank(t *testing.T) {
	s := []int{5, 4, 5, 5, 2, 1, 5, 6, 1, 3, 5, 0}
	w := NewWaveletMatrix(s)
	m := make(map[int]struct{})
	for _, v := range s {
		m[v] = struct{}{}
	}
	for k, _ := range m {
		for i := 1; i <= len(s); i++ {
			c := 0
			for j := 0; j < i; j++ {
				if s[j] == k {
					c++
				}
			}
			if a := w.Rank(k, i); c != a {
				t.Errorf("k == %v, i == %v", k, i)
				t.Errorf("%v != %v", c, a)
			}
		}
	}
}

func TestRankLess(t *testing.T) {
	s := []int{5, 4, 5, 5, 2, 1, 5, 6, 1, 3, 5, 0}
	w := NewWaveletMatrix(s)
	v := 6
	for l := 0; l < len(s); l++ {
		for r := l; r <= len(s); r++ {
			e := 0
			for i := l; i < r; i++ {
				if s[i] < v {
					e++
				}
			}
			if a := w.RankLess(l, r, v); e != a {
				t.Errorf("s == %v", s)
				t.Errorf("l == %v, r == %v, v == %v", l, r, v)
				t.Errorf("%v != %v", e, a)
				return
			}
		}
	}
}

func TestSelect(t *testing.T) {
	s := []int{5, 4, 5, 5, 2, 1, 5, 6, 1, 3, 5, 0}
	w := NewWaveletMatrix(s)
	m := make(map[int]struct{})
	for _, v := range s {
		m[v] = struct{}{}
	}
	for k, _ := range m {
		for i := 1; i <= len(s); i++ {
			e := 0
			c := 0
			for j := 0; j < len(s); j++ {
				if s[j] == k {
					c++
					if i == c {
						e = j + 1
					}
				}
			}
			if a := w.Select(k, i); e != a {
				t.Errorf("k == %v, i == %v", k, i)
				t.Errorf("%v != %v", e, a)
			}
		}
	}
}

func TestQuantile(t *testing.T) {
	s := []int{5, 4, 5, 5, 2, 1, 5, 6, 1, 3, 5, 0}
	w := NewWaveletMatrix(s)
	for l := 0; l < len(s); l++ {
		for r := l + 1; r <= len(s); r++ {
			e := make([]int, r-l)
			copy(e, s[l:r])
			sort.Ints(e)
			for k := 0; k < r-l; k++ {
				rank := k + 1
				if a := w.Quantile(l, r, rank); e[k] != a {
					t.Errorf("l == %v, r == %v, rank == %v", l, r, rank)
					t.Errorf("%v != %v", e[k], a)
				}
			}
		}
	}
}

// TestTopkOrder checks whether Topk return array is sort by frequency in descending order or not.
// The original order is not stable.
func TestTopkOrder(t *testing.T) {
	s := []int{5, 4, 5, 5, 2, 1, 5, 6, 1, 3, 5, 0}
	w := NewWaveletMatrix(s)
	for l := 0; l < len(s); l++ {
		for r := l + 1; r <= len(s); r++ {
			si := make([]int, r-l)
			copy(si, s[l:r])
			for k := 1; k <= len(si); k++ {
				a := w.Topk(l, r, k)
				f := a[0][1]
				for i, v := range a {
					nf := v[1]
					if f < nf {
						t.Errorf("The frequency order is incorrect.")
						t.Errorf("si == %v", si)
						t.Errorf("i == %v, si[i] == %v", i, si[i])
						t.Errorf("a == %v", a)
						break
					}
					f = nf
				}
			}
		}
	}
}

// TestTopkWithMaxk checks whether Topk returns correct value and frequency pairs or not.
func TestTopkWithMaxk(t *testing.T) {
	s := []int{5, 4, 5, 5, 2, 1, 5, 6, 1, 3, 5, 0}
	w := NewWaveletMatrix(s)
	for l := 0; l < len(s); l++ {
		for r := l + 1; r <= len(s); r++ {
			si := make([]int, r-l)
			copy(si, s[l:r])
			m := make(map[int]struct{})
			o := make([]int, 0) // value order appeared initially
			for _, v := range si {
				_, ok := m[v]
				if !ok {
					o = append(o, v)
					m[v] = struct{}{}
				}
			}
			e := make([][2]int, 0)
			for _, ov := range o {
				c := 0
				for _, tv := range si {
					if ov == tv {
						c++
					}
				}
				e = append(e, [2]int{ov, c})
			}
			sort.SliceStable(e, func(i, j int) bool { return e[i][1] > e[j][1] })
			k := len(e) // max k
			a := w.Topk(l, r, k)
			me := make(map[int]int)
			for _, v := range e {
				me[v[0]] = v[1]
			}
			ma := make(map[int]int)
			for _, v := range a {
				ma[v[0]] = v[1]
			}
			if !reflect.DeepEqual(me, ma) {
				t.Errorf("si == %v", si)
				t.Errorf("l == %v, r == %v, k == %v", l, r, k)
				t.Errorf("%v != %v", e, a)
				t.Errorf("%v != %v", me, ma)
			}
		}
	}
}

// TestTopkWithOverk checks whether Topk returns appropriate length slice with k is greater than the length of the slice or not.
func TestTopkWithOverk(t *testing.T) {
	s := []int{5, 4, 5, 5, 2, 1, 5, 6, 1, 3, 5, 0}
	w := NewWaveletMatrix(s)
	for l := 0; l < len(s); l++ {
		for r := l + 1; r <= len(s); r++ {
			si := make([]int, r-l)
			copy(si, s[l:r])
			m := make(map[int]struct{})
			o := make([]int, 0) // value order appeared initially
			for _, v := range si {
				_, ok := m[v]
				if !ok {
					o = append(o, v)
					m[v] = struct{}{}
				}
			}
			e := make([][2]int, 0)
			for _, ov := range o {
				c := 0
				for _, tv := range si {
					if ov == tv {
						c++
					}
				}
				e = append(e, [2]int{ov, c})
			}
			sort.SliceStable(e, func(i, j int) bool { return e[i][1] > e[j][1] })
			k := len(si) * 2 // over k
			a := w.Topk(l, r, k)
			me := make(map[int]int)
			for _, v := range e {
				me[v[0]] = v[1]
			}
			ma := make(map[int]int)
			for _, v := range a {
				ma[v[0]] = v[1]
			}
			if !reflect.DeepEqual(me, ma) {
				t.Errorf("si == %v", si)
				t.Errorf("l == %v, r == %v, k == %v", l, r, k)
				t.Errorf("%v != %v", e, a)
				t.Errorf("%v != %v", me, ma)
			}
		}
	}
}

// TestSum checks whether Sum returns summary of value in [l, r) or not
func TestSum(t *testing.T) {
	s := []int{5, 4, 5, 5, 2, 1, 5, 6, 1, 3, 5, 0}
	w := NewWaveletMatrix(s)
	for l := 0; l < len(s); l++ {
		for r := l + 1; r <= len(s); r++ {
			si := make([]int, r-l)
			copy(si, s[l:r])
			e := 0
			for _, v := range si {
				e += v
			}
			if a := w.Sum(l, r); e != a {
				t.Errorf("si == %v", si)
				t.Errorf("l == %v, r == %v", l, r)
				t.Errorf("%v != %v", e, a)
			}
		}
	}
}

// TestIntersect checks whether Intersect returns common values and its frequency in [l1, r1) and [l2, r2) or not.
func TestIntersect(t *testing.T) {
	s := []int{5, 4, 5, 5, 2, 1, 5, 6, 1, 3, 5, 0}
	w := NewWaveletMatrix(s)
	for l1 := 0; l1 < len(s); l1++ {
		for r1 := l1 + 1; r1 <= len(s); r1++ {
			for l2 := 0; l2 < len(s); l2++ {
				for r2 := l2 + 1; r2 <= len(s); r2++ {
					si1 := make([]int, r1-l1)
					copy(si1, s[l1:r1])
					m1 := make(map[int]int)
					for _, v := range si1 {
						m1[v]++
					}

					si2 := make([]int, r2-l2)
					copy(si2, s[l2:r2])
					m2 := make(map[int]int)
					for _, v := range si2 {
						m2[v]++
					}

					var e [][3]int
					for k, v1 := range m1 {
						v2, ok := m2[k]
						if ok {
							e = append(e, [3]int{k, v1, v2})
						}
					}

					a := w.Intersect(l1, r1, l2, r2)

					sort.Slice(e, func(i, j int) bool { return e[i][0] < e[j][0] })
					sort.Slice(a, func(i, j int) bool { return e[i][0] < e[j][0] })

					if !reflect.DeepEqual(e, a) {
						t.Errorf("si1 == %v", si1)
						t.Errorf("si2 == %v", si2)
						t.Errorf("l1 == %v, r1 == %v, l2 == %v, r2 == %v", l1, r1, l2, r2)
						t.Errorf("%T != %T", e, a)
						t.Errorf("%v != %v", e, a)
						t.Errorf("len %v != %v", len(e), len(a))
					}
				}
			}
		}
	}
}

// TestRangefreq checks whether Rangefreq returns number of values x <= v < y in [l, r) or not.
func TestRangefreq(t *testing.T) {
	s := []int{5, 4, 5, 5, 2, 1, 5, 6, 1, 3, 5, 0}
	w := NewWaveletMatrix(s)
	for l := 0; l < len(s); l++ {
		for r := l; r < len(s); r++ {
			for x := 0; x < 10; x++ {
				for y := x; y < 10; y++ {
					e := 0
					for i := l; i < r; i++ {
						if x <= s[i] && s[i] < y {
							e++
						}
					}
					a := w.Rangefreq(l, r, x, y)
					if e != a {
						t.Errorf("s == %v", s)
						t.Errorf("l == %v", l)
						t.Errorf("r == %v", r)
						t.Errorf("x == %v", x)
						t.Errorf("y == %v", y)
						t.Errorf("%v != %v", e, a)
					}
				}
			}
		}
	}
}

func BenchmarkNewWaveletMatrix(b *testing.B) {
	var s []int
	for i := 0; i < 200000; i++ {
		s = append(s, rand.Int())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewWaveletMatrix(s)
	}
}

func BenchmarkRangefreq(b *testing.B) {
	var s []int
	for i := 0; i < 200000; i++ {
		s = append(s, rand.Int())
	}
	w := NewWaveletMatrix(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := rand.Intn(len(s))
		if r == 0 {
			r++
		}
		l := rand.Intn(r)
		y := rand.Int()
		x := rand.Intn(y)
		w.Rangefreq(l, r, x, y)
	}
}
