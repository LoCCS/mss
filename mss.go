package mss

import (
	"io"

	"github.com/sammy00/mss/container/stack"
	"github.com/sammy00/mss/ots/winternitz"
)

func treeHash(H uint32, rand io.Reader) ([]byte, []*stack.Stack) {
	S := stack.New()

	// push the 1st leaf on stack
	sk, _ := winternitz.GenerateKey(rand)
	S.Push(&Node{0, winternitz.HashPk(&sk.PublicKey)})

	numLeaf := ((1 << H) - 1)
	for leaf := 1; leaf < numLeaf; leaf++ {
		sk, _ = winternitz.GenerateKey(rand)
		i, nu := 0, winternitz.HashPk(&sk.PublicKey)
		for !S.Empty() {
			node := S.Peek().(*Node)

			if node.height != i {
				break
			}

			S.Pop()
			i, nu = i+1, merge(node.nu, nu)
		}

		S.Push(&Node{i, nu})
	}

	root := S.Peek().(*Node)
	S.Pop()

	return root.nu, nil
}

// MerkleSS implements the Merkle signature scheme
type MerkleSS struct {
	H             uint32
	Auth          [][]byte
	retainedStack []*stack.Stack
	root          []byte
}

/*
// NewMerkleSS makes a fresh Merkle signing routine
//	by running the generate key and setup procedure
func NewMerkleSS(H uint32) (*MerkleSS, error) {
	if H < 2 {
		return nil, errors.New("H should be larger than 1")
	}

	S := stack.New()

	// push the 1st leaf on stack
	sk, _ := winternitz.GenerateKey(rand)
	S.Push(&Node{0, winternitz.HashPk(&sk.PublicKey)})

	numLeaf := ((1 << H) - 1)
	for leaf := 1; leaf < numLeaf; leaf++ {
		sk, _ = winternitz.GenerateKey(rand)
		i, nu := 0, winternitz.HashPk(&sk.PublicKey)
		for !S.Empty() {
			node := S.Peek().(*Node)

			if node.height != i {
				break
			}

			S.Pop()
			i, nu = i+1, merge(node.nu, nu)
		}

		S.Push(&Node{i, nu})
	}

	root := S.Peek()
}
*/
