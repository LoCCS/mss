package winternitz

import (
	"bytes"
	"testing"

	"github.com/LoCCS/mss/rand"
)

// TestKeyIteratorSerialize checks the serialization/deserialization
//	between KeyIterator and JSON
func TestKeyIteratorSerialization(t *testing.T) {
	seed, _ := rand.RandSeed()

	iter := NewKeyIterator(seed)
	iter.Next()

	integratedSeed := iter.Serialize()

	iter2 := new(KeyIterator)
	if iter2.Init(integratedSeed) {
		integratedSeed2 := iter2.Serialize()

		if !bytes.Equal(integratedSeed, integratedSeed2) {
			t.Fatal("error in MarshalAssJSON/UnmarshalFromJSON")
		}
	}
}

// TestKeyIteratorExec tests correctness KeyIterator, including
//	normal running (demo by iter) and recovered running from
//	seed (demo by iter2)
func TestKeyIteratorExec(t *testing.T) {
	seed, _ := rand.RandSeed()

	iter := NewKeyIterator(seed)
	iter.Next()

	iter2 := new(KeyIterator)
	if !iter2.Init(iter.Serialize()) {
		t.Fatal("invalid integrated seed")
	}

	for i := 0; i < 2; i++ {
		sk1, _ := iter.Next()
		sk2, _ := iter2.Next()

		if !IsEqual(sk1, sk2) {
			t.Fatal("sk's should be equal")
		}
	}
}
