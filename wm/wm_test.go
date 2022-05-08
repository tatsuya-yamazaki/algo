package wm

import (
	"sort"
	"testing"
)

func Test(t *testing.T) {
	t := []int{5,4,5,5,2,1,5,6,1,3,5,0}
	w := NewWaveletMatrix(s)
	m := make(map[int]struct{})
	for _, v := range t {
		m[v] = struct{}{}
	}
	for i:=0; i<len(s); i++ {
		if w.Access(i) != s[i] {
			t.Errorf("%v != %v", i, s[i])
		}
	}
}

func Test(t *testing.T) {
	t := []int{5,4,5,5,2,1,5,6,1,3,5,0}
	w := NewWaveletMatrix(s)
	m := make(map[int]struct{})
	for _, v := range t {
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
			if w.Rank(k,i) != c {
				t.Errorf("%v != %v", k, i)
			}
		}
	}
}

func Test(t *testing.T) {
	t := []int{5,4,5,5,2,1,5,6,1,3,5,0}
	w := NewWaveletMatrix(s)
	m := make(map[int]struct{})
	for _, v := range t {
		m[v] = struct{}{}
	}
	for k, _ := range m {
		for i:=0; i<len(s); i++ {
			a := len(s)
			c := 0
			for j:=0; j<len(s); j++ {
				if s[j] == k {
					c++
					if i == c - 1 {
						a = j
					}
				}
			}
			if w.Select(k, i) != a {
				t.Errorf("%v != %v", k, i)
			}
		}
	}
}

func Test(t *testing.T) {
	t := []int{5,4,5,5,2,1,5,6,1,3,5,0}
	w := NewWaveletMatrix(s)
	m := make(map[int]struct{})
	for _, v := range t {
		m[v] = struct{}{}
	}
	for l:=0; l<len(s); l++ {
		for r:=l+1; r<=len(s); r++ {
			for i:=0; i<r-l; i++ {
				qt := make([]int, r - l)
				copy(qt, s[l:r])
				sort.Ints(qt)
				if q := w.Quantile(l,r,i); q != qt[i] {
					t.Errorf("%v != %v", q, qt[i])
				}
			}
		}
	}
}
