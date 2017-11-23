package main

//import "github.com/sammy00/mss/container/stack"
import (
	"io"

	"github.com/sammy00/mss/container/stack"
	"github.com/sammy00/mss/ots/winternitz"
)

// Node is a node in the Merkle tree
type Node struct {
	height int
	index  int
	nu     []byte
}

func TreeHashImpr(H, K int, rand io.Reader) {
	S := stack.New()

	// push the 1st leaf on stack
	sk, _ := winternitz.GenerateKey(rand)
	S.Push(&Node{0, 0, winternitz.HashPk(&sk.PublicKey)})

	availableIndices := make([]int, K-1)
	leafIdxMax := ((1 << H) - 1)
	for s := 1; s <= leafIdxMax; s++ {

	}

	//S.Push(&Node{0, 0})
}
