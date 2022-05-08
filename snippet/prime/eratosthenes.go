package prime

import(
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	defer iou.Fl()

	n := iou.I()
	a := 0

	max := 1000001
	p := make([]bool, max)
	q := NewQueue()
	p[0] = true
	p[1] = true
	q.Push(2)
	pt := make([]int, 0)
	pt = append(pt, 2)
	for q.Next() {
		v := q.Pop()
		ac := v
		for ac < max {
			p[ac] = true
			ac += v
		}
		for i:=v+1; i<max; i++ {
			if !p[i] {
				q.Push(i)
				pt = append(pt, i)
				break
			}
		}
	}
}
