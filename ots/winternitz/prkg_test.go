package winternitz

import (
	"bytes"
	"testing"

	"github.com/LoCCS/mss/rand"
)

// TestKeyIteratorSerialize checks the serialization/deserialization
//	between KeyIterator and bytes sequence
func TestKeyIteratorSerialization(t *testing.T) {
	seed, _ := rand.RandSeed()

	iter := NewKeyIterator(seed)
	iter.Next()

	compositeSeed := iter.Serialize()

	iter2 := new(KeyIterator)
	if iter2.Init(compositeSeed) {
		compositeSeed2 := iter2.Serialize()

		if !bytes.Equal(compositeSeed, compositeSeed2) {
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

		// check equality
		for j := range sk1.x {
			if !bytes.Equal(sk1.x[j], sk2.x[j]) {
				t.Fatalf("want %x, got %x", sk1.x[j], sk2.x[j])
			}
		}
	}
}
