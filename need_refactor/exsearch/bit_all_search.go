package main
import (
	"fmt"
)

func main() {
	fmt.Println(BitAll(5))
}

func BitAll(n int) [][]int {
	ret := make([][]int, 1<<n)
	for i:=0; i<(1<<n); i++ {
		var r []int
		for j:=0; j<n; j++ {
			if i & (1<<j) > 0 {
				r = append(r, j)
			}
		}
		ret[i] = r
	}
	return ret
}
