package main
import (
	"fmt"
)

func main() {
	fmt.Println(bitAll(5))
}

func bitAll(n int) [][]int {
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
