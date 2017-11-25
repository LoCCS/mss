package mss

import (
	"errors"
	"fmt"
	"math"

	"github.com/sammy00/mss/container/stack"
	"github.com/sammy00/mss/ots/winternitz"
)

// Node is a node in the Merkle tree
type Node struct {
	height uint32
	//index  int
	nu []byte
}

// String generates the string representation of node
func (node *Node) String() string {
	return fmt.Sprintf("{height: %v}", node.height)
}

type TreeHashStack struct {
	leaf uint32
	//leafUpper uint32
	height    uint32
	nodeStack *stack.Stack
}

func NewTreeHashStack(startingLeaf, h uint32) *TreeHashStack {
	treeHashStack := new(TreeHashStack)
	treeHashStack.Init(startingLeaf, h)

	return treeHashStack
}
func (th *TreeHashStack) Init(startingLeaf, h uint32) error {
	if 1 != startingLeaf%2 {
		return errors.New("invalid index of starting leaf")
	}

	//th.leaf, th.leafUpper = startingLeaf, (1 << h)
	th.leaf, th.height = startingLeaf, h
	th.nodeStack = stack.New() // clear up the stack

	return nil
}

func (th *TreeHashStack) IsCompleted() bool {
	return th.leaf >= uint32(1<<th.height)
}

func (th *TreeHashStack) LowestTailHeight() uint32 {
	if th.nodeStack.Empty() {
		return th.height
	}
	if th.IsCompleted() {
		return math.MaxUint32
	}
	return th.Top().height
}
func (th *TreeHashStack) Top() *Node {
	if th.nodeStack.Empty() {
		return nil
	}

	return th.nodeStack.Peek().(*Node)
}

func (th *TreeHashStack) Update(numOp uint32, keyItr *winternitz.SkPkIterator) {
	if th.IsCompleted() {
		return
	}

	for (numOp > 0) && !th.IsCompleted() {
		if th.nodeStack.Len() < 2 {
			// invoke key generator to make a new leaf and
			//	add the new leaf to S
			sk, _ := keyItr.Next()
			th.nodeStack.Push(&Node{0, winternitz.HashPk(&sk.PublicKey)})
			th.leaf++
			numOp--
			continue
		}

		e1, e2 := th.nodeStack.Peek2()
		node1 := e1.(*Node)
		node2 := e2.(*Node)
		if node1.height == node2.height {
			th.nodeStack.Pop()
			th.nodeStack.Pop()

			th.nodeStack.Push(&Node{node1.height + 1, merge(node2.nu, node1.nu)})
		} else {
			// TODO: invoke key generator to make a new leaf and
			//	add the new leaf to S
			sk, _ := keyItr.Next()
			th.nodeStack.Push(&Node{0, winternitz.HashPk(&sk.PublicKey)})
			th.leaf++
		}

		numOp--
	}
}
