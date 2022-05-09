package prime

func Eratosthenes(max int) (p []int) {
	if max < 2 {
		return p
	}

	isNotPrime := make([]bool, max+1)
	isNotPrime[0] = true
	isNotPrime[1] = true

	for i:=2; i<=max; i++ {
		if isNotPrime[i] {
			continue
		}
		p = append(p, i)
		ac := i * 2
		for ac <= max {
			isNotPrime[ac] = true
			ac += i
		}
	}
	return p
}
