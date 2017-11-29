package winternitz

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/sammy00/mss/config"
)

// TestHashFuncApplier tests the HashFuncApplier
func TestHashFuncApplier(t *testing.T) {
	in := []byte("TestHashFuncApplier, good night")
	h := config.HashFunc()

	outs := make([][]byte, 16)
	for i := range outs {
		//fmt.Printf("*** numTimes=%v\n", i)

		applier := NewHashFuncApplier(big.NewInt(int64(i+1)), h)
		outs[i] = applier.Eval(in, nil)

		var prev []byte
		if i > 0 {
			prev = outs[i-1]
		} else {
			prev = in
		}

		h.Reset()
		h.Write(prev)
		if !bytes.Equal(h.Sum(nil), outs[i]) {
			t.Fatalf("want %s, but got %s",
				hex.EncodeToString(h.Sum(nil)), hex.EncodeToString(outs[i]))
		}
	}
}
