package container

import (
	"fmt"
	"testing"
)

// TestStack tests the correctness of the stack implementation
func TestStack(t *testing.T) {
	stack := NewStack()

	hash := []byte("HelloWorld, Stack-----------------------")
	for i := 0; i < 3; i++ {
		for j := 0; j < 7; j++ {
			hash[i*3+j] = byte('*')
			stack.Push(NewMerkleNode(i, j, hash))
		}
	}

	for !stack.Empty() {
		top, nextTop := stack.Peek2()
		fmt.Println("******")
		fmt.Println(top)
		fmt.Println(nextTop)
		stack.Pop()
	}
}
