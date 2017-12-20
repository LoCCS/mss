package winternitz

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"runtime"
	"sync"
)

// PublicKey as container for public key
type PublicKey struct {
	*WtnOpts
	Y [][]byte
}

// PrivateKey as container for private key,
// it also embeds its corresponding public key
type PrivateKey struct {
	PublicKey
	x [][]byte
}

// WinternitzSig as container for the Winternitz one-time signature
type WinternitzSig struct {
	sigma [][]byte
}

// Clone returns a copy of this public key
func (pk *PublicKey) Clone() *PublicKey {
	pkC := new(PublicKey)

	if nil != pk.WtnOpts {
		pkC.WtnOpts = pk.WtnOpts.Clone()
	}

	if nil != pk.Y {
		pkC.Y = make([][]byte, len(pk.Y))
		for i := range pk.Y {
			pkC.Y[i] = make([]byte, len(pk.Y[i]))
			copy(pkC.Y[i], pk.Y[i])
		}
	}

	return pkC
}

// GenerateKey generates a one-time key pair
// according to specification (nonce, key-pair-index) in opts and
// by getting randomness from rng
func GenerateKey(opts *WtnOpts, rng io.Reader) (*PrivateKey, error) {
	sk := new(PrivateKey)
	sk.x = make([][]byte, wtnLen)
	sk.Y = make([][]byte, wtnLen)

	// sample the private key
	skIdx := make([]byte, 32)
	seed := make([]byte, opts.SecurityLevel())
	rng.Read(seed)
	for i := range sk.x {
		sk.x[i] = make([]byte, opts.SecurityLevel())
		binary.BigEndian.PutUint32(skIdx, uint32(i))
		prf(sk.x[i], seed, skIdx)
	}
	sk.WtnOpts = opts.Clone()

	// evaluate the corresponding public key
	var wg sync.WaitGroup
	numCPU := uint32(runtime.NumCPU())
	workloadSize := (wtnLen + numCPU - 1) / numCPU
	for i := uint32(0); i < numCPU; i++ {
		wg.Add(1)

		// worker
		go func(workerIdx uint32) {
			defer wg.Done()

			// range of chains run by current worker
			from := workerIdx * workloadSize
			to := from + workloadSize
			if to > wtnLen {
				to = wtnLen
			}

			// make an address copy
			addr := make(address, sk.WtnOpts.addr.Len())
			copy(addr, sk.WtnOpts.addr)

			for j := from; j < to; j++ {
				// set the index of chain
				addr.setChainAddress(j)
				// derive the corresponding y[i]
				sk.Y[j] = evalChain(sk.x[j], 0, wtnMask, sk.WtnOpts.nonce, addr)
			}
		}(i)
	}

	// wait until all Y[j] has been computed
	wg.Wait()

	return sk, nil
}

// Sign generates the signature for a message digest based on
// the given private key
// the opts should have the same seed and key-pair index and
// as that of generating key pairs
func Sign(sk *PrivateKey, hash []byte) (*WinternitzSig, error) {
	blocks := hashToBlocks(hash)
	if len(sk.x) != len(blocks) {
		return nil, errors.New("mismatched secret key and b_i")
	}

	wtnSig := new(WinternitzSig)
	wtnSig.sigma = make([][]byte, len(sk.x))

	var wg sync.WaitGroup
	numCPU := uint32(7)
	xLen := uint32(len(sk.x))
	jobSize := (xLen + numCPU - 1) / numCPU
	for i := uint32(0); i < numCPU; i++ {
		wg.Add(1)
		// worker
		go func(workerIdx uint32) {
			defer wg.Done()
			// range of blocks evaluated by current worker
			from := workerIdx * jobSize
			to := from + jobSize
			if to > xLen {
				to = xLen
			}

			// make an address copy
			addr := make(address, sk.WtnOpts.addr.Len())
			copy(addr, sk.WtnOpts.addr)

			for j := from; j < to; j++ {
				// set index of chain
				addr.setChainAddress(j)
				// sigma_j=f^b_j(x_j)
				wtnSig.sigma[j] = evalChain(sk.x[j], 0, uint32(blocks[j]),
					sk.WtnOpts.nonce, addr)
			}
		}(i)
	}

	// synchronise all verification of Y[j]
	wg.Wait()

	return wtnSig, nil
}

// Verify verifies the Merkle signature on a message digest
// against the claimed public key and
// the opts should have the same seed and key-pair index
// as that of generating key pairs
func Verify(pk *PublicKey, hash []byte, wtnSig *WinternitzSig) bool {
	blocks := hashToBlocks(hash)
	if (len(pk.Y) != len(blocks)) || (len(pk.Y) != len(wtnSig.sigma)) {
		return false
	}

	// flag indicating signature is valid
	ok := true

	// verify the signature in parallel
	var wg sync.WaitGroup
	numCPU := uint32(runtime.NumCPU())
	yLen := uint32(len(pk.Y))
	jobSize := (yLen + numCPU - 1) / numCPU
	for i := uint32(0); i < numCPU; i++ {
		wg.Add(1)
		// worker
		go func(workerIdx uint32) {
			defer wg.Done()
			// range of blocks evaluated by current worker
			from := workerIdx * jobSize
			to := from + jobSize
			if to > yLen {
				to = yLen
			}

			// make an address copy
			addr := make(address, pk.WtnOpts.addr.Len())
			copy(addr, pk.WtnOpts.addr)

			for j := from; j < to; j++ {
				// set index of chain
				addr.setChainAddress(j)
				// f^{2^w-1-b_j}(sigma_j)
				y := evalChain(wtnSig.sigma[j], uint32(blocks[j]),
					wtnMask-uint32(blocks[j]), pk.WtnOpts.nonce, addr)

				ok = ok && bytes.Equal(pk.Y[j], y)
			}
		}(i)
	}

	// synchronise all verification of Y[j]
	wg.Wait()

	return ok
}

//Serialize encodes the winternitz signature
func (sig *WinternitzSig) Serialize() []byte{
	sNum := len(sig.sigma)
	size := 0
	if sNum > 0 && sig.sigma[0] != nil {
		size = len(sig.sigma[0])
	}
	sigBytes := make([]byte, 2 + 2 + sNum * size)
	binary.LittleEndian.PutUint16(sigBytes[0:2], uint16(sNum))
	binary.LittleEndian.PutUint16(sigBytes[2:4], uint16(size))
	offset := 4
	for _, s := range sig.sigma{
		copy(sigBytes[offset: offset + size], s)
		offset += size
	}
	return sigBytes
}

//Deserialize decodes the winternitz signature from bytes
func DeserializeWinternitzSig(sigBytes []byte) *WinternitzSig{
	sNum := int(binary.LittleEndian.Uint16(sigBytes[0:2]))
	size := int(binary.LittleEndian.Uint16(sigBytes[2:4]))
	offset := 4
	sigma := make([][]byte, sNum)
	for i := 0; i < int(sNum); i++{
		sigma[i] = sigBytes[offset : offset + size]
		offset += size
	}
	return &WinternitzSig{
		sigma,
	}
}
