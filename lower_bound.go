package main

import (
	"fmt"
)

func main() {
	s := []int{1,2,2,3,3,3,4,4,4,4,5,5,5,5,5}
	fmt.Println(LowerBound(s, 0))
	fmt.Println(LowerBound(s, 1))
	fmt.Println(LowerBound(s, 2))
	fmt.Println(LowerBound(s, 3))
	fmt.Println(LowerBound(s, 4))
	fmt.Println(LowerBound(s, 5))
	fmt.Println(LowerBound(s, 6))
	fmt.Println(LowerBound(s, 7))
	fmt.Println(LowerBound(s, 8))
}

func LowerBound(s []int, value int) int {
	l, r := 0, len(s)
	for l != r {
		m := (l + r) / 2
		if value > s[m] {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}
