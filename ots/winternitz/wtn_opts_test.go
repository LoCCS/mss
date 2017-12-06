package winternitz

import (
	"bytes"
	"testing"
)

func TestWtnOptsCopy(t *testing.T) {
	opts := NewWtnOpts(SecurityLevel)
	optsC := opts.Clone()

	optsC.nonce[0] = optsC.nonce[0] ^ 0xff
	optsC.SetKeyIdx(0x1234)

	if bytes.Equal([]byte(opts.addr), []byte(optsC.addr)) {
		t.Fatal("modification after copy, address should not be equal")
	}

	if bytes.Equal(opts.nonce, optsC.nonce) {
		t.Fatal("modification after copy, nonce should not be equal")
	}
}
