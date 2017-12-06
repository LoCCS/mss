package winternitz

import (
	mathrand "math/rand"
	"testing"

	"github.com/LoCCS/mss/config"
	mssrand "github.com/LoCCS/mss/rand"
)

// TestWinternitzSig tests the signing/verifying of W-OTS
func TestWinternitzSig(t *testing.T) {
	hashFunc := config.HashFunc()
	hashFunc.Write([]byte("hello Merkle signature scheme..."))
	// compute digest
	hash := hashFunc.Sum(nil)

	// generate keys
	sk, _ := GenerateKey(DummyWtnOpts, mssrand.Reader)

	wtnSig, err := Sign(sk, hash)
	if nil != err {
		t.Fatal(err)
	}

	if !Verify(&sk.PublicKey, hash, wtnSig) {
		t.Fatal("verification failed")
	}
}

func TestIterativeWinternitzSig(t *testing.T) {
	hashFunc := config.HashFunc()
	hashFunc.Write([]byte("hello Merkle signature scheme..."))
	// compute digest
	hash := hashFunc.Sum(nil)

	seed0, _ := mssrand.RandSeed()
	keyItr := NewKeyIterator(seed0)
	for i := uint32(0); i < 32; i++ {
		sk, _ := keyItr.Next()

		wtnSig, err := Sign(sk, hash)
		if nil != err {
			t.Fatal(err)
		}

		if !Verify(&sk.PublicKey, hash, wtnSig) {
			t.Fatal("verification failed")
		}
	}
}

// TestWinternitzSig tests the signing/verifying of W-OTS
//	with corrupted public key
func TestWinternitzSigBadPk(t *testing.T) {
	hashFunc := config.HashFunc()
	hashFunc.Write([]byte("Testing Winternitz One-Time Signature..."))
	// compute digest
	hash := hashFunc.Sum(nil)

	// generate keys
	sk, _ := GenerateKey(DummyWtnOpts, mssrand.Reader)

	wtnSig, err := Sign(sk, hash)
	if nil != err {
		t.Fatal(err)
	}

	pk := &sk.PublicKey
	// corrupt some byte of pk
	i, j := mathrand.Int()%len(pk.Y), mathrand.Int()%len(pk.Y[0])
	pk.Y[i][j] ^= 0xff

	if Verify(pk, hash, wtnSig) {
		t.Fatal("verification failed")
	}
}
