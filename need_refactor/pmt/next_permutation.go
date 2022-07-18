package main

import "fmt"

func main() {
	s := []int{0, 1, 1, 2, 2, 2, 3}
	c := 0
	for {
		c++
		fmt.Println(s)
		if NextPermutation(s) {
			break
		}
	}
	fmt.Println(c)
}

func NextPermutation(s []int) bool {
	l, b := -1, -1
	for i := 0; i < len(s)-1; i++ {
		if s[i] < s[i+1] {
			l = i
		}
	}
	if l == -1 {
		return true
	}
	for i := len(s) - 1; i >= 0; i-- {
		if s[l] < s[i] {
			b = i
			break
		}
	}
	s[l], s[b] = s[b], s[l]
	for i := 1; i <= (len(s)-1-l)/2; i++ {
		s[l+i], s[len(s)-i] = s[len(s)-i], s[l+i]
	}
	return false
}
