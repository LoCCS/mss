package winternitz

import (
	"testing"

	"github.com/LoCCS/mss/config"
	mssrand "github.com/LoCCS/mss/rand"
)

func BenchmarkGenerateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateKey(DummyWtnOpts, mssrand.Reader)
	}
}

func BenchmarkWOTS(b *testing.B) {
	sha := config.HashFunc()
	sha.Write([]byte("hello Merkle signature scheme..."))
	// compute digest
	hash := sha.Sum(nil)

	for i := 0; i < b.N; i++ {
		sk, err := GenerateKey(DummyWtnOpts, mssrand.Reader)

		wtnSig, err := Sign(sk, hash)
		if nil != err {
			b.Fatal(err)
		}

		if !Verify(&sk.PublicKey, hash, wtnSig) {
			b.Fatal("verification failed")
		}
	}
}
