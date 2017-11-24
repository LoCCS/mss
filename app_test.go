package mss_test

import (
	"fmt"
	"testing"
)

func TestApp(t *testing.T) {
	i := uint32(3)
	fmt.Println(1<<i - 1)
}
