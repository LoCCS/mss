package winternitz

import (
	"hash"
	"math/big"

	"github.com/sammy00/mss/config"
)

// HashFuncApplier composes a composite function `f(x)=h^(numTimes)(x)`
//	based on a given hash function `h`
type HashFuncApplier struct {
	numTimes *big.Int
	h        hash.Hash
}

// NewHashFuncApplier allocates and returns a new HashFuncApplier
//	based on the given `numTimes` and primitive hash function `h`
func NewHashFuncApplier(numTimes *big.Int, h hash.Hash) *HashFuncApplier {
	return &HashFuncApplier{numTimes, h}
}

// Eval applies the underlying primitive hash function to the given
//	input `in` iteratively `numTimes` times
func (applier *HashFuncApplier) Eval(in []byte, numTimes *big.Int) []byte {
	delta := big.NewInt(1)

	numItr := new(big.Int)
	if nil != numTimes {
		numItr.Set(numTimes)
	} else {
		numItr.Set(applier.numTimes)
	}

	out := in
	for ; numItr.Sign() > 0; numItr.Sub(numItr, delta) {
		// update `out` as `out=h(out)`
		applier.h.Reset()
		applier.h.Write(out)
		out = applier.h.Sum(nil)
	}

	return out
}

// HashPk computes the hash value for a MSS public key
func HashPk(pk *PublicKey) []byte {
	hashFunc := config.HashFunc()

	for i := range pk.Y {
		hashFunc.Write(pk.Y[i])
	}

	return hashFunc.Sum(nil)
}
