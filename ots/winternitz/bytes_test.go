package winternitz

import (
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
