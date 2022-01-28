package main

import(
	"fmt"
	"math"
)

func main() {
	fmt.Println(primeFactorize(138600))
	fmt.Println(primeFactorize(1000000007))
}

func primeFactorize(n int) (ret [][2]int) {
	root := int(math.Sqrt(float64(n)))
	for i:=2; i<=root; i++ {
		var factor [2]int
		for n % i == 0 {
			n /= i
			factor[1]++
		}
		if factor[1] != 0 {
			factor[0] = i
			ret = append(ret, factor)
		}
	}
	if n != 1 {
		ret = append(ret, [2]int{n, 1})
	}
	return
}
