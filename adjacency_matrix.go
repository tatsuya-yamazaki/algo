package main

import(
	"fmt"
)

func main() {
	nam := NewAdjacencyMatrix(8)
	fmt.Println(nam)
}

func NewAdjacencyMatrix(n int) [][]bool {
        m := make([][]bool, n)
        for i:=0; i<n; i++ {
                m[i] = make([]bool, n)
        }
        return m
}
