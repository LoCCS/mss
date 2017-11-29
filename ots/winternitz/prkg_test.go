package winternitz

import (
	"bytes"
	"testing"

	"github.com/sammy00/mss/rand"
)

// TestKeyIteratorCoding checks the encoding/decoding
//	between KeyIterator and JSON
func TestKeyIteratorCoding(t *testing.T) {
	seed, _ := rand.RandSeed()

	iter := NewKeyIterator(seed)
	iter.Next()

	bytesJson, err := iter.MarshalAsJSON()
	if nil != err {
		t.Fatal(err)
	}

	iter2 := new(KeyIterator)
	iter2.UnmarshalFromJSON(bytesJson)
	bytesJson2, err := iter2.MarshalAsJSON()

	if !bytes.Equal(bytesJson, bytesJson2) {
		t.Fatal("error in MarshalAssJSON/UnmarshalFromJSON")
	}
}

// TestKeyIteratorExec tests correctness KeyIterator, including
//	normal running (demo by iter) and recovered running from
//	seed (demo by iter2)
func TestKeyIteratorExec(t *testing.T) {
	seed, _ := rand.RandSeed()

	iter := NewKeyIterator(seed)
	iter.Next()

	//iter2 := NewKeyIterator(iter.rng.ExportSeed())
	bytesJson, err := iter.MarshalAsJSON()
	if nil != err {
		t.Fatal(err)
	}
	iter2 := new(KeyIterator)
	iter2.UnmarshalFromJSON(bytesJson)

	for i := 0; i < 2; i++ {
		sk1, _ := iter.Next()
		sk2, _ := iter2.Next()

		if !IsEqual(sk1, sk2) {
			t.Fatal("sk's should be equal")
		}
	}
}
