// TODO randomize test data. create function generateRandomData.
package wm

import (
	"reflect"
	"sort"
	"testing"
)

func TestAccess(t *testing.T) {
	s := []int{5,4,5,5,2,1,5,6,1,3,5,0}
	w := NewWaveletMatrix(s)
	for i:=0; i<len(s); i++ {
		if a := w.Access(i); s[i] != a {
			t.Errorf("%v != %v", s[i], a)
		}
	}
}

func TestRank(t *testing.T) {
	s := []int{5,4,5,5,2,1,5,6,1,3,5,0}
	w := NewWaveletMatrix(s)
	m := make(map[int]struct{})
	for _, v := range s {
		m[v] = struct{}{}
	}
	for k, _ := range m {
		for i:=0; i<len(s); i++ {
			c := 0
			for j:=0; j<=i; j++ {
				if s[j] == k {
					c++
				}
			}
			if a := w.Rank(k,i); c != a {
				t.Errorf("k == %v, i == %v", k, i)
				t.Errorf("%v != %v", c, a)
			}
		}
	}
}

func TestSelect(t *testing.T) {
	s := []int{5,4,5,5,2,1,5,6,1,3,5,0}
	w := NewWaveletMatrix(s)
	m := make(map[int]struct{})
	for _, v := range s {
		m[v] = struct{}{}
	}
	for k, _ := range m {
		for i:=0; i<len(s); i++ {
			e := len(s)
			c := 0
			for j:=0; j<len(s); j++ {
				if s[j] == k {
					c++
					if i == c - 1 {
						e = j
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
	s := []int{5,4,5,5,2,1,5,6,1,3,5,0}
	w := NewWaveletMatrix(s)
	for l:=0; l<len(s); l++ {
		for r:=l+1; r<=len(s); r++ {
			e := make([]int, r - l)
			copy(e, s[l:r])
			sort.Ints(e)
			for k:=0; k<r-l; k++ {
				if a := w.Quantile(l,r,k); e[k] != a {
					t.Errorf("l == %v, r == %v, k == %v", l, r, k)
					t.Errorf("%v != %v", e[k], a)
				}
			}
		}
	}
}

// TestTopkOrder checks whether Topk return array is sort by frequency in descending order or not.
// The original order is not stable.
func TestTopkOrder(t *testing.T) {
	s := []int{5,4,5,5,2,1,5,6,1,3,5,0}
	w := NewWaveletMatrix(s)
	for l:=0; l<len(s); l++ {
		for r:=l+1; r<=len(s); r++ {
			si := make([]int, r - l)
			copy(si, s[l:r])
			for k:=1; k<=len(si); k++ {
				a := w.Topk(l,r,k)
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
	s := []int{5,4,5,5,2,1,5,6,1,3,5,0}
	w := NewWaveletMatrix(s)
	for l:=0; l<len(s); l++ {
		for r:=l+1; r<=len(s); r++ {
			si := make([]int, r - l)
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
			a := w.Topk(l,r,k)
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
	s := []int{5,4,5,5,2,1,5,6,1,3,5,0}
	w := NewWaveletMatrix(s)
	for l:=0; l<len(s); l++ {
		for r:=l+1; r<=len(s); r++ {
			si := make([]int, r - l)
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
			a := w.Topk(l,r,k)
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
	s := []int{5,4,5,5,2,1,5,6,1,3,5,0}
	w := NewWaveletMatrix(s)
	for l:=0; l<len(s); l++ {
		for r:=l+1; r<=len(s); r++ {
			si := make([]int, r - l)
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
	s := []int{5,4,5,5,2,1,5,6,1,3,5,0}
	w := NewWaveletMatrix(s)
	for l1:=0; l1<len(s); l1++ {
		for r1:=l1+1; r1<=len(s); r1++ {
			for l2:=0; l2<len(s); l2++ {
				for r2:=l2+1; r2<=len(s); r2++ {
					si1 := make([]int, r1 - l1)
					copy(si1, s[l1:r1])
					m1 := make(map[int]int)
					for _, v := range si1 {
						m1[v]++
					}

					si2 := make([]int, r2 - l2)
					copy(si2, s[l2:r2])
					m2 := make(map[int]int)
					for _, v := range si2 {
						m2[v]++
					}

					e := make([][3]int, 0)
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
						t.Errorf("%v != %v", e, a)
					}
				}
			}
		}
	}
}

