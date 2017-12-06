package winternitz

import (
	"bytes"
	"errors"
	"io"
)

// PublicKey as container for public key
type PublicKey struct {
	*WtnOpts
	Y [][]byte
}

// PrivateKey as container for private key,
//	it also embeds its corresponding public key
type PrivateKey struct {
	PublicKey
	x [][]byte
}

// WinternitzSig as container for the Winternitz one-time signature
type WinternitzSig struct {
	sigma [][]byte
}

// GenerateKey generates a one-time key pair
//	according to specification (nonce, key-pair-index) in opts and
//	by getting randomness from rng
func GenerateKey(opts *WtnOpts, rng io.Reader) (*PrivateKey, error) {
	sk := new(PrivateKey)
	sk.x = make([][]byte, wtnLen)
	sk.Y = make([][]byte, wtnLen)

	// sample the private key
	for i := range sk.x {
		sk.x[i] = make([]byte, opts.SecurityLevel())
		rng.Read(sk.x[i])
	}
	sk.WtnOpts = opts.Clone()

	// evaluate the corresponding public key
	numIter := uint32((1 << w) - 1)
	for i := uint32(0); i < wtnLen; i++ {
		// set the index of chain
		opts.addr.setChainAddress(i)
		// derive the corresponding y[i]
		sk.Y[i] = evalChain(sk.x[i], 0, numIter, opts)
	}

	return sk, nil
}

// Sign generates the signature for a message digest based on
//	the given private key
//	the opts should have the same seed and key-pair index and
//	as that of generating key pairs
//func Sign(opts *WtnOpts, sk *PrivateKey, hash []byte) (*WinternitzSig, error) {
func Sign(sk *PrivateKey, hash []byte) (*WinternitzSig, error) {
	blocks := hashToBlocks(hash)
	if len(sk.x) != len(blocks) {
		return nil, errors.New("mismatched secret key and b_i")
	}

	wtnSig := new(WinternitzSig)
	wtnSig.sigma = make([][]byte, len(sk.x))

	opts := sk.WtnOpts.Clone()
	for i := range sk.x {
		// set index of chain
		opts.addr.setChainAddress(uint32(i))
		// sigma_i=f^b_i(x_i)
		wtnSig.sigma[i] = evalChain(sk.x[i], 0, uint32(blocks[i]), opts)
	}

	return wtnSig, nil
}

// Verify verifies the Merkle signature on a message digest
//	against the claimed public key and
//	the opts should have the same seed and key-pair index
//	as that of generating key pairs
//func Verify(opts *WtnOpts, pk *PublicKey, hash []byte, wtnSig *WinternitzSig) bool {
func Verify(pk *PublicKey, hash []byte, wtnSig *WinternitzSig) bool {
	blocks := hashToBlocks(hash)
	if (len(pk.Y) != len(blocks)) || (len(pk.Y) != len(wtnSig.sigma)) {
		return false
	}

	// w-1
	wBaseMax := uint32((1 << w) - 1)

	opts := pk.Clone()
	for i := range wtnSig.sigma {
		opts.addr.setChainAddress(uint32(i))
		// f^{w-1-b_i}(sigma_i)
		y := evalChain(wtnSig.sigma[i], uint32(blocks[i]), wBaseMax-uint32(blocks[i]), opts)

		if !bytes.Equal(pk.Y[i], y) {
			return false
		}
	}

	return true
}
