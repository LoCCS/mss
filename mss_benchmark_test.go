package mss

import (
	"testing"

	"github.com/LoCCS/mss/rand"
)

func BenchmarkMSSSetup(b *testing.B) {
	const H = 16
	seed, _ := rand.RandSeed()

	for i := 0; i < b.N; i++ {
		if _, err := NewMerkleAgent(H, seed); nil != err {
			b.Fatalf("unexpected error in  NewMerkleAgent(%v,%x)", H, seed)
		}
	}
}

func BenchmarkMSSStd(b *testing.B) {
	const H = 16
	seed, _ := rand.RandSeed()
	merkleAgent, err := NewMerkleAgent(H, seed)
	if nil != err {
		b.Fatal("unexpected error in setting up")
	}

	b.ResetTimer()
	// make a random message hash
	msg, _ := rand.RandSeed()
	// what if no more leaf to use in the Merkle agent
	for i := 0; i < b.N; i++ {
		_, sig, err := Sign(merkleAgent, msg)
		if nil != err {
			b.Fatalf("error in signing %x", msg)
		}

		if !Verify(merkleAgent.Root(), msg, sig) {
			b.Fatal("verification failed")
		}
	}
}
