package winternitz

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"
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
