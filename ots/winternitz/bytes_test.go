package winternitz

import (
	"bytes"
	"testing"

	"golang.org/x/crypto/sha3"
)

/*
func TestGetUint64(t *testing.T) {
	buf := []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12}
	want := uint64(0x123456789abcdef0)

	x := GetUint64(buf)
	if x != want {
		t.Fatalf("want %x, got %x", want, x)
	}
}
*/

func TestToBase(t *testing.T) {
	const baseWidth = 4
	X := []byte{0x12, 0x34}

	outLen := 8 * len(X) / w
	out := make([]byte, outLen)

	ToBase(out, X, w)

	want := []byte{1, 2, 3, 4}
	if len(want) != len(out) {
		t.Fatal("error output length")
	}
	if !bytes.Equal(want, out) {
		t.Fatalf("wants %x, got %x", want, out)
	}
}

/*
// TestToBytes checks big-endian serialization by ToBytes()
func TestToBytes(t *testing.T) {
	x := uint64(0x123456789abcdef0)

	xBytesBE := make([]byte, 8)
	ToBytes(xBytesBE, x)
	xWants := []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}
	if !bytes.Equal(xBytesBE, xWants) {
		t.Fatal("invalid output: wants %x, got %x", xWants, xBytesBE)
	}

	y := uint64(0x1234567)

	yBytesBE := make([]byte, 4)
	ToBytes(yBytesBE, y)
	yWants := []byte{0x01, 0x23, 0x45, 0x67}
	if !bytes.Equal(yBytesBE, yWants) {
		t.Fatal("invalid output: wants %x, got %x", yWants, yBytesBE)
	}
}

// TestPutBytesAndGetUint64 ensures GetUint64() is the inverse to ToBytes()
func TestPutBytesAndGetUint64(t *testing.T) {
	x := rand.Uint64()

	y := make([]byte, 8)
	ToBytes(y, x)

	z := GetUint64(y)

	if x != z {
		t.Fatalf("want %x, got %x", x, z)
	}
}
*/

// TestHashToBlocks checks the encoding of hash byte slice
//	to blocks
func TestHashToBlocks(t *testing.T) {
	hash := sha3.Sum256([]byte("TestHashToBlocks"))
	blocks := hashToBlocks(hash[:])

	wm := byte((1 << w) - 1)
	for _, b := range blocks {
		if b > wm {
			t.Fatalf("%v>%v", b, wm)
		}
	}
}
