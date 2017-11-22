package main

import (
	"fmt"
	"testing"
)

func hello(world []byte) []byte {
	world2 := make([]byte, len(world))
	copy(world2, world)
	return world2
}
func TestApp(t *testing.T) {
	world := []byte("world")
	world2 := hello(world)

	world[0] = byte(77)
	world2[1] = byte(88)

	fmt.Println(world)
	fmt.Println(world2)
}
