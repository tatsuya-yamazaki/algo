package main

import(
	"fmt"
	"math"
)

func main() {
	s := make([]int, 0)

	s = append(s, 3)
	s = append(s, 5)
	s = append(s, 2)
	s = append(s, 11)
	s = append(s, 9)
	s = append(s, 1)
	s = append(s, 20)
	s = append(s, 8)

	q := NewRangeMaximumQuery(s)

	fmt.Println(q.s)
	fmt.Println(q.Query(3, 5))
}

type RangeMaximumQuery struct {
        s []int
        start int
}

func NewRangeMaximumQuery(seq []int) *RangeMaximumQuery {
        n := len(seq)
        t := 1
        for t < n {
                t *= 2
        }
        ret := &RangeMaximumQuery{}
        ret.s = make([]int, t * 2 - 1)
        ret.start = t - 1
	for i:=0; i<len(ret.s); i++ {
		ret.s[i] = math.MinInt64
        }
        for i, v := range seq {
                ret.Set(i, v)
        }
        return ret
}

func (q *RangeMaximumQuery) Set(i, v int) {
	i = q.start + i
	q.s[i] = v
	for i > 0 {
		i = i / 2 + i % 2 - 1
		if q.s[i*2+1] > q.s[i*2+2] {
			q.s[i] = q.s[i*2+1]
		} else {
			q.s[i] = q.s[i*2+2]
		}
	}
}

func (q *RangeMaximumQuery) Query(a, b int) int {
	return q.query(a, b+1, 0, 0, len(q.s) - q.start)
}

func (q *RangeMaximumQuery) query(a, b, i, l, r int) int {
	if r <= a || b <= l {
		return math.MinInt64
	} else if a <= l && r <= b {
		return q.s[i]
	} else {
		vl := q.query(a, b, i*2+1, l, (l+r)/2)
		vr := q.query(a, b, i*2+2, (l+r)/2, r)
		if vl > vr {
			return vl
		} else {
			return vr
		}
	}
}
