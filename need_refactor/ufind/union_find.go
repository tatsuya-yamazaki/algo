package main

import (
	"fmt"
)

func main() {
	uf := NewUnionFind(5)

	uf.Unite(1, 3)
	uf.Unite(2, 4)
	uf.Unite(4, 0)

	fmt.Println("0,1", uf.SameRoot(0, 1))
	fmt.Println("0,2", uf.SameRoot(0, 2))
	fmt.Println("0,3", uf.SameRoot(0, 3))
	fmt.Println("0,4", uf.SameRoot(0, 4))
	fmt.Println("1,2", uf.SameRoot(1, 2))
	fmt.Println("1,3", uf.SameRoot(1, 3))
	fmt.Println("1,4", uf.SameRoot(1, 4))
	fmt.Println("2,3", uf.SameRoot(2, 3))
	fmt.Println("2,4", uf.SameRoot(2, 4))
	fmt.Println("3,4", uf.SameRoot(3, 4))
	fmt.Println(uf.parent)
	fmt.Println(uf.rank)
}

type UnionFind struct {
	parent []int
	rank   []int
}

func NewUnionFind(length int) *UnionFind {
	parent := make([]int, length)
	rank := make([]int, length)
	for i := 0; i < length; i++ {
		parent[i] = i
	}
	return &UnionFind{parent, rank}
}

func (u *UnionFind) Root(index int) int {
	if u.parent[index] == index {
		return index
	} else {
		u.parent[index] = u.Root(u.parent[index])
		return u.parent[index]
	}
}

func (u *UnionFind) SameRoot(a, b int) bool {
	return u.Root(a) == u.Root(b)
}

func (u *UnionFind) Unite(a, b int) {
	a = u.Root(a)
	b = u.Root(b)

	if a == b {
		return
	}

	if u.rank[a] < u.rank[b] {
		u.parent[a] = b
	} else {
		u.parent[b] = a
		if u.rank[a] == u.rank[b] {
			u.rank[a]++
		}
	}
}
