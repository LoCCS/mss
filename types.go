package mss

import (
	"errors"
	"fmt"

	"github.com/sammy00/mss/container/stack"
)

// Node is a node in the Merkle tree
type Node struct {
	height int
	//index  int
	nu []byte
}

// String generates the string representation of node
func (node *Node) String() string {
	return fmt.Sprintf("{height: %v}", node.height)
}

type TreeHashStack struct {
	leaf      uint32
	leafUpper uint32
	S         *stack.Stack
}

func (th *TreeHashStack) Init(startingLeaf, h uint32) error {
	if 1 != startingLeaf%2 {
		return errors.New("invalid index of starting leaf")
	}

	th.leaf, th.leafUpper = startingLeaf, (1 << h)
	th.S = stack.New() // clear up the stack

	return nil
}

func (th *TreeHashStack) IsCompleted() bool {
	return th.leaf >= th.leafUpper
}

func (th *TreeHashStack) Update(numOp uint32) {
	if th.IsCompleted() {
		return
	}

	for (numOp > 0) && !th.IsCompleted() {
		if th.S.Len() < 2 {
			// invoke key generator to make a new leaf and
			//	add the new leaf to S
			continue
		}

		e1, e2 := th.S.Peek2()
		node1 := e1.(*Node)
		node2 := e2.(*Node)
		if node1.height == node2.height {
			th.S.Pop()
			th.S.Pop()

			th.S.Push(&Node{node1.height + 1, merge(node2.nu, node1.nu)})
		} else {
			// TODO: invoke key generator to make a new leaf and
			//	add the new leaf to S
		}

		numOp--
	}
}
