package winternitz

import (
	"bytes"
	"errors"
	"io"

	"github.com/sammy00/mss/config"
)

// PublicKey as container for public key
type PublicKey struct {
	Y [][]byte
}

// PrivateKey as container for private key,
//	it also embeds its corresponding public key
type PrivateKey struct {
	PublicKey
	x [][]byte
}

// WinternitzSig as container for the Merkle signature
type WinternitzSig struct {
	sigma [][]byte
}

// GenerateKey generates a one-time key pair
func GenerateKey(rand io.Reader) (*PrivateKey, error) {
	sk := new(PrivateKey)
	sk.x = make([][]byte, t)
	sk.Y = make([][]byte, t)

	numIter := uint32(w - 1)
	for i := range sk.x {
		sk.x[i] = make([]byte, config.Size)
		// make a rand x[i]
		rand.Read(sk.x[i])

		// derive the corresponding y[i]
		sk.Y[i] = HashChainEval(sk.x[i], numIter)
	}

	return sk, nil
}

// Sign generates the signature for a message digest based on
//	the given private key
func Sign(sk *PrivateKey, hash []byte) (*WinternitzSig, error) {
	blocks := hashToBlocks(hash)
	if len(sk.x) != len(blocks) {
		return nil, errors.New("mismatched secret key and b_i")
	}

	wtnSig := new(WinternitzSig)
	wtnSig.sigma = make([][]byte, len(sk.x))

	for i := range sk.x {
		wtnSig.sigma[i] = HashChainEval(sk.x[i], uint32(blocks[i]))
	}

	return wtnSig, nil
}

// Verify verifies the Merkle signature on a message digest
//	against the claimed public key
func Verify(pk *PublicKey, hash []byte, wtnSig *WinternitzSig) bool {
	blocks := hashToBlocks(hash)
	if (len(pk.Y) != len(blocks)) || (len(pk.Y) != len(wtnSig.sigma)) {
		return false
	}

	// w-1
	numIter := uint32(w - 1)
	for i := range wtnSig.sigma {
		// f^{w-1-b_i}(sigma[i])
		y := HashChainEval(wtnSig.sigma[i], numIter-uint32(blocks[i]))

		if !bytes.Equal(pk.Y[i], y) {
			return false
		}
	}

	return true
}
