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
