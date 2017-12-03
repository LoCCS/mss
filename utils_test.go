package mss

import (
	"bytes"
	"testing"

	"github.com/LoCCS/mss/config"
)

func TestMerge(t *testing.T) {
	hashLeft, hashRight := []byte("hello"), []byte("world")

	hashMerged := merge(hashLeft, hashRight)

	hashConcatenated := append(hashLeft, hashRight...)
	sha := config.HashFunc()
	sha.Write(hashConcatenated)
	hashMerged2 := sha.Sum(nil)

	if !bytes.Equal(hashMerged, hashMerged2) {
		t.Fatal("failed")
	}
}
