package main
import(
	"fmt"
)

type BinarySearchTree struct {
	root *Node
}

func NewBinarySearchTree () *BinarySearchTree {
	return &BinarySearchTree{}
}

type Node struct {
	value int
	right *Node
	left *Node
}

func NewNode(value int) *Node {
	return &Node{value: value}
}

func (n *Node) Value() int {
	return n.value
}

func (t *BinarySearchTree) Find(value int) (n *Node, route []*Node) {
	n = t.root
	for n != nil {
		route = append(route, n)
		if value == n.value {
			return n, route
		} else if value > n.value {
			n = n.right
		} else {
			n = n.left
		}
	}
	return nil, route
}

func (t *BinarySearchTree) Add(value int) *Node {
	if t.root == nil {
		t.root = NewNode(value)
		return t.root
	}
	_, route := t.Find(value)
	parent := route[len(route)-1]
	if parent.value == value {
		return parent
	}
	nn := NewNode(value)
	if value > parent.value {
		parent.right = nn
	} else {
		parent.left = nn
	}
	return parent
}

func (t *BinarySearchTree) Remove(value int) bool {
	if t.root == nil {
		return false
	}
	n, route := t.Find(value)
	if n == nil {
		return false
	}
	parent := t.getParentFromRoute(route)
	if n.left != nil && n.right != nil {
		t.removeNodeHasBothChildren(n)
	} else if n.left != nil {
		t.removeNodeHasLeft(n, parent)
	} else if n.right != nil {
		t.removeNodeHasRight(n, parent)
	} else {
		t.removeNodeHasNoChild(n, parent)
	}
	return true
}

func (*BinarySearchTree) getParentFromRoute(route []*Node) *Node {
	if len(route) > 1 {
		return route[len(route)-2]
	}
	return nil
}

func (t *BinarySearchTree) removeNodeHasBothChildren(n *Node) {
	leftMax, route := t.Max(n.left)
	leftMaxParent := n
	if leftMax != n.left {
		leftMaxParent = t.getParentFromRoute(route)
	}
	n.value = leftMax.value
	if leftMax.left == nil {
		t.removeNodeHasNoChild(leftMax, leftMaxParent)
	} else {
		t.removeNodeHasLeft(leftMax, leftMaxParent)
	}
}

func (t *BinarySearchTree) removeNodeHasLeft(n, parent *Node) {
	t.replaceNode(n, parent, n.left)
	n.left = nil
}

func (t *BinarySearchTree) removeNodeHasRight(n, parent *Node) {
	t.replaceNode(n, parent, n.right)
	n.right = nil
}

func (t *BinarySearchTree) removeNodeHasNoChild(n, parent *Node) {
	t.replaceNode(n, parent, nil)
}

func (t *BinarySearchTree) replaceNode(n, parent, newNode *Node) {
	if parent != nil {
		if parent.left == n {
			parent.left = newNode
		} else {
			parent.right = newNode
		}
	}
	if n == t.root {
		t.root = newNode
	}
}

func (t *BinarySearchTree) Max(n *Node) (max *Node, route []*Node) {
	if t.root == nil {
		return nil, nil
	}
	if n == nil {
		n = t.root
	}
	route = append(route, n)
	for n.right != nil {
		n = n.right
		route = append(route, n)
	}
	return n, route
}

func (t *BinarySearchTree) Min(n *Node) (min *Node, route []*Node) {
	if t.root == nil {
		return nil, nil
	}
	if n == nil {
		n = t.root
	}
	route = append(route, n)
	for n.left != nil {
		n = n.left
		route = append(route, n)
	}
	return n, route
}

func (t *BinarySearchTree) Echo() {
	if t.root == nil {
		fmt.Println("nil")
		return
	}
	t.echo(t.root, "")
}

func (t *BinarySearchTree) echo(n *Node, space string) {
	space += "    "
	if n.right != nil {
		t.echo(n.right, space)
	}
	fmt.Println(space, n.value)
	if n.left != nil {
		t.echo(n.left, space)
	}
}

func (t *BinarySearchTree) Repl() {
	for {
		s := ""
		n := 0
		fmt.Scan(&s)
		switch s {
		case "a":
			fmt.Scan(&n)
			t.Add(n)
			t.Echo()
		case "r":
			fmt.Scan(&n)
			t.Remove(n)
			t.Echo()
		case "f":
			fmt.Scan(&n)
			fmt.Println(t.Find(n))
		case "p":
			t.Echo()
		}
	}
}

func main() {
	bst := NewBinarySearchTree()
	bst.Repl()
}
