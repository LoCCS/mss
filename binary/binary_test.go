package binary

import (
	"bytes"
	"testing"
)

func TestPutUint(t *testing.T) {
	var v uint32 = 0x12345678

	b := make([]byte, 3)
	want := []byte{0x34, 0x56, 0x78}

	PutUint(b, uint64(v))
	if !bytes.Equal(b, want) {
		t.Logf("want %x, got %x", want, b)
	}
}
