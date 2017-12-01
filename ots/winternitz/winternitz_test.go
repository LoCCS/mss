package winternitz

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mathrand "math/rand"
	"testing"

	mrand "github.com/sammy00/mss/rand"
)

// TestGenerateKey tests the generation of a one-time key pair (sk,pk)
func TestGenerateKey(t *testing.T) {
	fmt.Println("***1st generation")
	sk, _ := GenerateKey(rand.Reader)

	fmt.Println("totally", len(sk.x), "key pairs as")
	fmt.Println("{")
	for i, x := range sk.x {
		fmt.Printf(" (%s,\n  %s)\n", hex.EncodeToString(x), hex.EncodeToString(sk.Y[i]))
	}
	fmt.Println("}")

	fmt.Println("***2nd generation")
	sk, _ = GenerateKey(rand.Reader)

	fmt.Println("totally", len(sk.x), "key pairs as")
	fmt.Println("{")
	for i, x := range sk.x {
		fmt.Printf(" (%s,\n  %s)\n", hex.EncodeToString(x), hex.EncodeToString(sk.Y[i]))
	}
	fmt.Println("}")
}

// TestWinternitzSig tests the signing/verifying of W-OTS
func TestWinternitzSig(t *testing.T) {
	hashFunc := HashFunc()
	hashFunc.Write([]byte("hello Merkle signature scheme..."))
	// compute digest
	hash := hashFunc.Sum(nil)

	// generate keys
	sk, _ := GenerateKey(mrand.Reader)
	//fmt.Println(len(sk.x))

	wtnSig, err := Sign(sk, hash)
	if nil != err {
		t.Fatal(err)
	}

	if !Verify(&sk.PublicKey, hash, wtnSig) {
		t.Fatal("verification failed")
	}
}

// TestWinternitzSig tests the signing/verifying of W-OTS
//	with corrupted public key
func TestWinternitzSigBadPk(t *testing.T) {
	hashFunc := HashFunc()
	hashFunc.Write([]byte("hello Merkle signature scheme..."))
	// compute digest
	hash := hashFunc.Sum(nil)

	// generate keys
	sk, _ := GenerateKey(mrand.Reader)
	//fmt.Println(len(sk.x))

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
