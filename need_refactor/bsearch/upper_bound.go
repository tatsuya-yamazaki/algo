package main

import (
	"fmt"
)

func main() {
	s := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 5}
	fmt.Println(UpperBound(s, 0))
	fmt.Println(UpperBound(s, 1))
	fmt.Println(UpperBound(s, 2))
	fmt.Println(UpperBound(s, 3))
	fmt.Println(UpperBound(s, 4))
	fmt.Println(UpperBound(s, 5))
	fmt.Println(UpperBound(s, 6))
	fmt.Println(UpperBound(s, 7))
	fmt.Println(UpperBound(s, 8))
}

func UpperBound(s []int, value int) int {
	l, r := 0, len(s)
	for l != r {
		m := (l + r) / 2
		if value < s[m] {
			r = m
		} else {
			l = m + 1
		}
	}
	return l
}
