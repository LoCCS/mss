package main

import (
	"fmt"
	"testing"
)

func TestApp(t *testing.T) {
	hello := 5
	fmt.Printf("%v\n", hello)

	world := hello--
	fmt.Printf("%v\n", world)
}
