package winternitz

import (
	"crypto/rand"
	"testing"
)

func BenchmarkEvalChain(b *testing.B) {
	hash := make([]byte, SecurityLevel)
	for i := 0; i < b.N; i++ {
		rand.Read(hash)
		evalChain(hash, 0, wtnMask, DummyWtnOpts.nonce, DummyWtnOpts.addr)
	}
}
