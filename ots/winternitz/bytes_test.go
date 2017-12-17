package winternitz

import (
	"bytes"
	"testing"

	"golang.org/x/crypto/sha3"
)

func TestToBase(t *testing.T) {
	const baseWidth = 4
	X := []byte{0x12, 0x34}

	outLen := 8 * len(X) / baseWidth
	out := make([]byte, outLen)

	ToBase(out, X, baseWidth)

	want := []byte{1, 2, 3, 4}
	if len(want) != len(out) {
		t.Fatal("error output length")
	}
	if !bytes.Equal(want, out) {
		t.Fatalf("wants %x, got %x", want, out)
	}
}

// TestHashToBlocks checks the encoding of hash byte slice
// to blocks
func TestHashToBlocks(t *testing.T) {
	hash := sha3.Sum256([]byte("TestHashToBlocks"))
	blocks := hashToBlocks(hash[:])

	for _, b := range blocks {
		if b > wtnMask {
			t.Fatalf("%v>%v", b, wtnMask)
		}
	}
}
