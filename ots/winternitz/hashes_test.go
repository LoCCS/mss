package winternitz

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// TestHashFSum tests keyed hash func calculating
//	(0x00||key||data)
func TestFSum(t *testing.T) {
	hash := make([]byte, SecurityLevel)
	key := make([]byte, SecurityLevel)
	msg := []byte("Hello World")

	fSum(hash, key, msg)

	want, _ := hex.DecodeString("8f33f161f7ba43cabc1292e06623a30dce2f69cbed4a346c78c679b70917cc5b")
	if !bytes.Equal(hash, want) {
		t.Fatalf("want %x, got %x\n", want, hash)
	}
}

// TestPrf tests pseudo random function calculating
//	(0x03||key||data)
func TestPrf(t *testing.T) {
	hash := make([]byte, SecurityLevel)
	key := make([]byte, SecurityLevel)
	msg := []byte("Hello World")

	prf(hash, key, msg)

	want, _ := hex.DecodeString("6eaf22c1d9e08e57b0f41de26a246d2337c46db13d082e4a62c152c3c05f1dda")
	if !bytes.Equal(hash, want) {
		t.Fatalf("want %x, got %x\n", want, hash)
	}
}
