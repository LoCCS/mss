package container

import (
	"fmt"
)

// MerkleNode is a node of a merkle tree
type MerkleNode struct {
	height int
	index  int
	nu     []byte
}

// NewMerkleNode makes a node with specific height, index
//	and hash digest
func NewMerkleNode(h, i int, hash []byte) *MerkleNode {
	node := &MerkleNode{height: h, index: i}

	if nil != hash {
		node.nu = make([]byte, len(hash))
		copy(node.nu, hash)
	}

	return node
}

// String returns the string representation of a node
func (node *MerkleNode) String() string {
	return fmt.Sprintf("{height: %v, index: %v, nu: %s}",
		node.height, node.index, string(node.nu))
}

// Element is an element in a single linked list
type Element struct {
	Value *MerkleNode
	next  *Element
}

// Next accesses the next element connected to this current node
func (e *Element) Next() *Element {
	return e.next
}

// Stack represents a stack
// The zero value for Stack is an unintialized stack
// ought to initialize by invoking Init()
type Stack struct {
	top  *Element // top element of the stack
	size int      // size of the stack
}

// New returns an initialized stack
func NewStack() *Stack {
	return &Stack{nil, 0}
}

// Push adds an element to the stack
func (s *Stack) Push(x *MerkleNode) {
	s.top = &Element{x, s.top}
	s.size++
}

// Pop out the top element from the stack
func (s *Stack) Pop() *MerkleNode {
	var x *MerkleNode
	if nil != s.top {
		x, s.top = s.top.Value, s.top.next
		s.size--
	}

	return x
}

// Peek returns the element in the top of the s
func (s *Stack) Peek() *MerkleNode {
	if nil != s.top {
		return s.top.Value
	}

	return nil
}

// Peek returns the two elements in the top of the s
func (s *Stack) Peek2() (*MerkleNode, *MerkleNode) {
	var x, y *MerkleNode

	if nil != s.top {
		x = s.top.Value

		if nil != s.top.next {
			y = s.top.next.Value
		}
	}

	return x, y
}

// Len returns the size of the s
func (s *Stack) Len() int {
	return s.size
}

// Empty returns true if the stack is empty
func (s *Stack) Empty() bool {
	return (0 == s.size)
}
