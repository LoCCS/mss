package stack

import (
	"fmt"
	"testing"
)

// TestStack tests the correctness of the stack implementation
func TestStack(t *testing.T) {
	stack := New()

	for i := 0; i < 8; i++ {
		stack.Push(i)
	}

	for !stack.Empty() {
		top, nextTop := stack.Peek2()
		fmt.Println("******")

		nextTopValue := -1
		if nil != nextTop {
			nextTopValue = nextTop.(int)
		}
		fmt.Println(top, nextTopValue)
		stack.Pop()
	}
}
