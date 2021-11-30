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
	parent *Node
	right *Node
	left *Node
}

func NewNode(value int) *Node {
	return &Node{value: value}
}

func (t *BinarySearchTree) Find(value int) (n *Node) {
	n = t.root
	for n != nil {
		if value == n.value {
			return n
		} else if value > n.value {
			n = n.right
		} else if value < n.value {
			n = n.left
		}
	}
	return nil
}

func (t *BinarySearchTree) Add(value int) *Node {
	if t.root == nil {
		t.root = NewNode(value)
		return t.root
	}
	n := t.root
	nn := NewNode(value)
	for n != nil {
		if value >= n.value {
			if n.right == nil {
				n.right = nn
				break
			} else {
				n = n.right
			}
		} else {
			if n.left == nil {
				n.left = nn
				break
			} else {
				n = n.left
			}
		}
	}
	nn.parent = n
	return nn
}

func (t *BinarySearchTree) Remove(value int) bool {
	if t.root == nil {
		return false
	}
	n := t.Find(value)
	if n == nil {
		return false
	}
	if n.left != nil && n.right != nil {
		t.removeNodeHasBothChildren(n)
	} else if n.left != nil {
		t.removeNodeHasLeft(n)
	} else if n.right != nil {
		t.removeNodeHasRight(n)
	} else {
		t.removeNodeHasNoChild(n)
	}
	return true
}

func (t *BinarySearchTree) removeNodeHasBothChildren(n *Node) {
	leftMax := t.Max(n.left)
	n.value = leftMax.value
	if leftMax.left == nil {
		t.removeNodeHasNoChild(leftMax)
	} else {
		t.removeNodeHasLeft(leftMax)
	}
}

func (t *BinarySearchTree) removeNodeHasLeft(n *Node) {
	t.replaceNode(n, n.left)
	n.left = nil
}

func (t *BinarySearchTree) removeNodeHasRight(n *Node) {
	t.replaceNode(n, n.right)
	n.right = nil
}

func (t *BinarySearchTree) removeNodeHasNoChild(n *Node) {
	t.replaceNode(n, nil)
}

func (t *BinarySearchTree) replaceNode(n, newNode *Node) {
	if n.parent != nil {
		if n.parent.left == n {
			n.parent.left = newNode
		} else if n.parent.right == n {
			n.parent.right = newNode
		}
		if newNode != nil {
			newNode.parent = n.parent
		}
		n.parent = nil
	}
	if n == t.root {
		t.root = newNode
	}
}

func (t *BinarySearchTree) Max(n *Node) *Node {
	if t.root == nil {
		return nil
	}
	if n == nil {
		n = t.root
	}
	for n.right != nil {
		n = n.right
	}
	return n
}

func (t *BinarySearchTree) Min(n *Node) *Node {
	if t.root == nil {
		return nil
	}
	if n == nil {
		n = t.root
	}
	for n.left != nil {
		n = n.left
	}
	return n
}

func (t *BinarySearchTree) Echo() {
	if t.root == nil {
		fmt.Println("nil")
		return
	}
	t.echo(t.root, "")
}

func (t *BinarySearchTree) echo(n *Node, space string) {
	if n.right != nil {
		t.echo(n.right, "  "+space)
	}
	fmt.Println(space, n.value)
	if n.left != nil {
		t.echo(n.left, "  "+space)
	}
}

func main() {
	bst := NewBinarySearchTree()
	for {
		s := ""
		n := 0
		fmt.Scan(&s)
		switch s {
		case "a":
			fmt.Scan(&n)
			bst.Add(n)
			bst.Echo()
		case "r":
			fmt.Scan(&n)
			bst.Remove(n)
			bst.Echo()
		case "f":
			fmt.Scan(&n)
			fmt.Println(bst.Find(n))
		case "p":
			bst.Echo()
		}
	}
}
