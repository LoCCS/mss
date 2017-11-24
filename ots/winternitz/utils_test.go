package winternitz

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/sammy00/mss/config"
	"github.com/sammy00/mss/rand"
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

// TestSkPkIterator tests correctness SkPkIterator, including
//	normal running (demo by iter) and recovered running from
//	seed (demo by iter2)
func TestSkPkIterator(t *testing.T) {
	seed, _ := rand.RandSeed()

	iter := NewSkPkIterator(seed)
	iter.Next()

	iter2 := NewSkPkIterator(iter.Seed())
	for i := 0; i < 2; i++ {
		sk1, _ := iter.Next()
		sk2, _ := iter2.Next()

		if !IsEqual(sk1, sk2) {
			t.Fatal("sk's should be equal")
		}
	}
}
