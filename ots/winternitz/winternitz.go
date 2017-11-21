package winternitz

import (
	"io"

	"github.com/sammy00/mss/config"
)

const (
	W = 5
)

// PublicKey as container for public key
type PublicKey struct {
	y [][]byte
}

// PrivateKey as container for private key,
//	it also embeds its corresponding public key
type PrivateKey struct {
	PublicKey
	x [][]byte
}

// GenerateKey generates a one-time key pair
func GenerateKey(rand io.Reader) (*PrivateKey, error) {
	sk := new(PrivateKey)
	sk.x = make([][]byte, t)
	sk.y = make([][]byte, len(sk.x))

	applier := NewHashFuncApplier(pow2ToW(), config.HashFunc())
	for i := range sk.x {
		sk.x[i] = make([]byte, config.Size)
		// make a rand x[i]
		rand.Read(sk.x[i])

		// derive the corresponding y[i]
		sk.y[i] = applier.Eval(sk.x[i])
	}

	return sk, nil
}

/*
func Sign(sk *PrivateKey, hash []byte) (*mssSig, error) {
	blocks, err := parseBlocks(hash)
	if nil != err {
		return nil, err
	}

	return nil, nil
}
*/
