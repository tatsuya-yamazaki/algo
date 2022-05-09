package prime

import (
	"testing"
	"reflect"
)

func TestPrimeFactorize(t *testing.T) {
	e1 := [][2]int{[2]int{2, 3}, [2]int{3, 2}, [2]int{5, 2}, [2]int{7, 1}, [2]int{11, 1}}
	a1 := PrimeFactorize(138600)
	if !reflect.DeepEqual(e1, a1) {
		t.Errorf("%v !=  %v", e1, a1)
		t.Errorf("%T !=  %T", e1, a1)
	}
	e2 := [][2]int{[2]int{1000000007, 1}}
	a2 := PrimeFactorize(1000000007)
	if !reflect.DeepEqual(e2, a2) {
		t.Errorf("%v !=  %v", e2, a2)
		t.Errorf("%T !=  %T", e2, a2)
	}
}

func TestEratosthenes(t *testing.T) {
	p := Eratosthenes(10000)
	for _, a := range p {
		if e := PrimeFactorize(a); len(e) != 1 {
			t.Errorf("a is %v", a)
			t.Errorf("prime factorize a is %v", e)
		}
	}
}
