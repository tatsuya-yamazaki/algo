package main

import (
	"fmt"
)

func main() {
	s := []int{1,2,2,3,3,3,4,4,4,4,5,5,5,5,5}
	fmt.Println(lowerBound(s, 0))
	fmt.Println(lowerBound(s, 1))
	fmt.Println(lowerBound(s, 2))
	fmt.Println(lowerBound(s, 3))
	fmt.Println(lowerBound(s, 4))
	fmt.Println(lowerBound(s, 5))
	fmt.Println(lowerBound(s, 6))
	fmt.Println(lowerBound(s, 7))
	fmt.Println(lowerBound(s, 8))
}

func lowerBound(s []int, value int) int {
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
