package winternitz

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// TestAddress checks basic functions, including
//	+ setKeyIdx
//	+ setChainAddress
//	+ setHashAddress
//	+ onMask
//	+ onKey
func TestAddress(t *testing.T) {
	addr := newAddress()

	wantInHex := []string{
		"0000000000000000000000000000000050607080000000000000000000000000",
		"0000000000000000000000000000000050607080102030400000000000000000",
		"0000000000000000000000000000000050607080102030401234567800000000",
		"00000000000000000000000000000000506070801020304012345678ffffffff",
		"0000000000000000000000000000000050607080102030401234567800000000",
	}

	addr.setKeyIdx(0x50607080)
	want, _ := hex.DecodeString(wantInHex[0])
	if !bytes.Equal([]byte(addr), want) {
		t.Fatalf("want %x, got\n", want, []byte(addr))
	}

	addr.setChainAddress(0x10203040)
	want, _ = hex.DecodeString(wantInHex[1])
	if !bytes.Equal([]byte(addr), want) {
		t.Fatalf("want %x, got\n", want, []byte(addr))
	}

	addr.setHashAddress(0x12345678)
	want, _ = hex.DecodeString(wantInHex[2])
	if !bytes.Equal([]byte(addr), want) {
		t.Fatalf("want %x, got\n", want, []byte(addr))
	}

	addr.onMask()
	want, _ = hex.DecodeString(wantInHex[3])
	if !bytes.Equal([]byte(addr), want) {
		t.Fatalf("want %x, got\n", want, []byte(addr))
	}

	addr.onKey()
	want, _ = hex.DecodeString(wantInHex[4])
	if !bytes.Equal([]byte(addr), want) {
		t.Fatalf("want %x, got\n", want, []byte(addr))
	}
}

func TestAddressCopy(t *testing.T) {
	addr := newAddress()
	addr.setKeyIdx(0x1234)
	addr.setChainAddress(0x5678)
	addr.setHashAddress(0x9abc)
	addr.onMask()

	addrC := make(address, 32)
	copy(addrC, addr)

	if !bytes.Equal(addr, addrC) {
		t.Fatalf("want %x, got %x", addr, addrC)
	}
}
