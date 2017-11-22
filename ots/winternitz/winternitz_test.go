package winternitz

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/sammy00/mss/config"
	mrand "github.com/sammy00/mss/rand"
)

// TestGenerateKey tests the generation of a one-time key pair (sk,pk)
func TestGenerateKey(t *testing.T) {
	sk, _ := GenerateKey(rand.Reader)

	fmt.Println("totally", len(sk.x), "key pairs as")
	fmt.Println("{")
	for i, x := range sk.x {
		fmt.Printf(" (%s,\n  %s)\n", hex.EncodeToString(x), hex.EncodeToString(sk.y[i]))
	}
	fmt.Println("}")
}

func TestMSS(t *testing.T) {
	hashFunc := config.HashFunc()
	hashFunc.Write([]byte("hello Merkle signature scheme..."))
	// compute digest
	hash := hashFunc.Sum(nil)
	// derive private keys
	sk, _ := GenerateKey(mrand.Reader)

	merkleSig, err := Sign(sk, hash)
	if nil != err {
		t.Fatal(err)
	}

	/*
		fmt.Println("signature is {")
		for _, sigma := range merkleSig.sigma {
			fmt.Printf(" %s,\n", hex.EncodeToString(sigma))
		}
		fmt.Println("}")
	*/

	if !Verify(&sk.PublicKey, hash, merkleSig) {
		t.Fatal("verification failed")
	}
}

func Test0MSSDebug(t *testing.T) {
	hashFunc := config.HashFunc()
	hashFunc.Write([]byte("hello Merkle signature scheme..."))
	// compute digest
	hash := hashFunc.Sum(nil)
	// derive private keys
	sk, _ := GenerateKey(mrand.Reader)

	merkleSig, err := Sign(sk, hash)
	if nil != err {
		t.Fatal(err)
	}

	//x := sk.x[0]
	//y := sk.y[0]

	/*
		fmt.Println("signature is {")
		for _, sigma := range merkleSig.sigma {
			fmt.Printf(" %s,\n", hex.EncodeToString(sigma))
		}
		fmt.Println("}")
	*/

	if !Verify(&sk.PublicKey, hash, merkleSig) {
		t.Fatal("verification failed")
	}
}
