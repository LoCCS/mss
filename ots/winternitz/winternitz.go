package winternitz

import (
	"bytes"
	"errors"
	"io"
	"math/big"

	"github.com/sammy00/mss/config"
)

const (
	W = 16 // number of bits to manipulate simultaneously
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

	applier := NewHashFuncApplier(bitMask(), config.HashFunc())

	for i := range sk.x {
		sk.x[i] = make([]byte, config.Size)
		// make a rand x[i]
		rand.Read(sk.x[i])

		// derive the corresponding y[i]
		sk.Y[i] = applier.Eval(sk.x[i], nil)
	}

	return sk, nil
}

// Sign generates the signature for a message digest based on
//	the given private key
func Sign(sk *PrivateKey, hash []byte) (*WinternitzSig, error) {
	blocks, err := hashToBlocks(hash)
	if nil != err {
		return nil, err
	}
	if len(sk.x) != len(blocks) {
		return nil, errors.New("mismatched private key and b_i")
	}

	merkleSig := new(WinternitzSig)
	merkleSig.sigma = make([][]byte, len(sk.x))

	applier := NewHashFuncApplier(nil, config.HashFunc())
	for i, x := range sk.x {
		merkleSig.sigma[i] = applier.Eval(x, blocks[i])
	}

	return merkleSig, nil
}

// Verify verifies the Merkle signature on a message digest
//	against the claimed public key
func Verify(pk *PublicKey, hash []byte, merkleSig *WinternitzSig) bool {
	blocks, err := hashToBlocks(hash)
	if (nil != err) || (len(pk.Y) != len(blocks)) ||
		(len(pk.Y) != len(merkleSig.sigma)) {
		return false
	}

	applier := NewHashFuncApplier(nil, config.HashFunc())
	mask := bitMask()
	numTimes := new(big.Int)
	for i := range merkleSig.sigma {
		// 2^W-1-b_i
		numTimes.Sub(mask, blocks[i])
		y := applier.Eval(merkleSig.sigma[i], numTimes)
		if !bytes.Equal(pk.Y[i], y) {
			return false
		}
	}

	return true
}
