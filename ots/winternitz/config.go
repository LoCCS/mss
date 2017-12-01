package winternitz

import (
	"hash"
	"math"

	"golang.org/x/crypto/sha3"
)

const (
	w = 16 // the Winternitz parameter, it is a member of the set {4, 16}
)

// length in bytes of digest produced by the employed hash function
var hashSize int

// parameters derived from w
// t is the number of n-byte string elements in a WOTS+ private key, public key, and signature.
var len1, len2, t uint32

func init() {
	hashSize = HashFunc().Size()

	wBits := float64(math.Ilogb(float64(w)))
	len1 = uint32(math.Ceil(float64(hashSize) * 8 / wBits))
	len2 = uint32(math.Floor(math.Log2(float64(len1)*(w-1))/wBits)) + 1
	t = len1 + len2
}

// returns the hash function to use
func HashFunc() hash.Hash {
	return sha3.New256()
}
