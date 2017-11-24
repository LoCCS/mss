package winternitz

import (
	"bytes"
	"hash"
	"math/big"

	"github.com/sammy00/mss/config"
	"github.com/sammy00/mss/rand"
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

// SkPkIterator is a iterator producing a key chain for
//	user based on a seed
type SkPkIterator struct {
	rng *rand.Rand
}

// NewSkPkIterator makes a key pair iterator
func NewSkPkIterator(seed []byte) *SkPkIterator {
	return &SkPkIterator{rand.New(seed)}
}

// Next estimates and returns the next sk-pk pair
func (iter *SkPkIterator) Next() (*PrivateKey, error) {
	return GenerateKey(iter.rng)
}

// Seed returns the internal updated seed for usage
//	such as saving state of the iterator
func (iter *SkPkIterator) Seed() []byte {
	return iter.rng.TellMeSeed()
}

// HashPk computes the hash value for a MSS public key
func HashPk(pk *PublicKey) []byte {
	hashFunc := config.HashFunc()

	for i := range pk.Y {
		hashFunc.Write(pk.Y[i])
	}

	return hashFunc.Sum(nil)
}

func IsEqual(sk1, sk2 *PrivateKey) bool {
	for i := range sk1.x {
		if !bytes.Equal(sk1.x[i], sk2.x[i]) {
			return false
		}
	}

	return true
}
