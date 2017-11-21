package winternitz

import (
	"hash"
	"math/big"
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
//	input `in` iteratively a predefined number of time
func (applier *HashFuncApplier) Eval(in []byte) []byte {
	numItr := new(big.Int)
	delta := big.NewInt(1)

	out := in
	for numItr.Set(applier.numTimes); numItr.Sign() > 0; numItr.Sub(numItr, delta) {
		// update `out` as `out=h(out)`
		applier.h.Reset()
		applier.h.Write(out)
		out = applier.h.Sum(nil)
	}

	return out
}
