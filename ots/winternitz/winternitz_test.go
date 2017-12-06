package winternitz

import (
	mathrand "math/rand"
	"testing"

	mssrand "github.com/LoCCS/mss/rand"
)

// TestWinternitzSig tests the signing/verifying of W-OTS
func TestWinternitzSig(t *testing.T) {
	hashFunc := HashFunc()
	hashFunc.Write([]byte("hello Merkle signature scheme..."))
	// compute digest
	hash := hashFunc.Sum(nil)

	// generate keys
	sk, _ := GenerateKey(DummyWtnOpts, mssrand.Reader)

	wtnSig, err := Sign(DummyWtnOpts, sk, hash)
	if nil != err {
		t.Fatal(err)
	}

	if !Verify(DummyWtnOpts, &sk.PublicKey, hash, wtnSig) {
		t.Fatal("verification failed")
	}
}

// TestWinternitzSig tests the signing/verifying of W-OTS
//	with corrupted public key
func TestWinternitzSigBadPk(t *testing.T) {
	hashFunc := HashFunc()
	hashFunc.Write([]byte("Testing Winternitz One-Time Signature..."))
	// compute digest
	hash := hashFunc.Sum(nil)

	// generate keys
	sk, _ := GenerateKey(DummyWtnOpts, mssrand.Reader)

	wtnSig, err := Sign(DummyWtnOpts, sk, hash)
	if nil != err {
		t.Fatal(err)
	}

	pk := &sk.PublicKey
	// corrupt some byte of pk
	i, j := mathrand.Int()%len(pk.Y), mathrand.Int()%len(pk.Y[0])
	pk.Y[i][j] ^= 0xff

	if Verify(DummyWtnOpts, pk, hash, wtnSig) {
		t.Fatal("verification failed")
	}
}
