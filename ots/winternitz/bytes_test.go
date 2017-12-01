package winternitz

import (
	"bytes"
	"fmt"
	"math"
	"testing"
)

// TestToBaseW checks the correctness of ToBaseW()
func TestToBaseW(t *testing.T) {
	const w = 16
	X := []byte{0x12, 0x34}

	outLen := 8 * len(X) / math.Ilogb(w)
	out := make([]byte, outLen)

	ToBaseW(out, X, w)

	want := []byte{1, 2, 3, 4}
	if len(want) != len(out) {
		t.Fatal("error output length")
	}
	for i := range out {
		if want[i] != out[i] {
			t.Fatalf("error output byte: wants %v, got %v", want[i], out[i])
		}
	}
}

// TestToBytes checks big-endian serialization by ToBytes()
func TestToBytes(t *testing.T) {
	X := uint64(0x123456789abcdef0)

	XBytesBE := ToBytes(X)
	XWants := []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}
	if !bytes.Equal(XBytesBE, XWants) {
		t.Fatal("invalid output: wants %x, got %x", XWants, XBytesBE)
	}

	Y := uint64(0x123456789abcde)

	YBytesBE := ToBytes(Y)
	YWants := []byte{0x00, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde}
	if !bytes.Equal(YBytesBE, YWants) {
		t.Fatal("invalid output: wants %x, got %x", YWants, YBytesBE)
	}
}

// TestHashToBlocks checks the encoding of hash byte slice
//	to blocks
func TestHashToBlocks(t *testing.T) {
	msg := "TestHashToBlocksV21"
	sha := HashFunc()
	sha.Write([]byte(msg))
	hash := sha.Sum(nil)

	blocks := hashToBlocks(hash)
	fmt.Printf("%x\n", blocks)
}
