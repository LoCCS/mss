package winternitz

import (
	"hash"
	"math"

	"golang.org/x/crypto/sha3"
)

// security parameter
const (
	SecurityLevel256 = 32 // 256 bits
	SecurityLevel512 = 64 // 512 bits
	SecurityLevel    = SecurityLevel256
)

// length in bytes of digest produced by the employed hash function
var hashSize int

const (
	w = 16 // the Winternitz parameter, it is a member of the set {4, 16}
)

// parameters derived from w
// wtnLen is the number of n-byte string elements in a WOTS+ private key, public key, and signature.
var wtnLen1, wtnLen2, wtnLen uint32

func init() {
	hashSize = HashFunc().Size()

	wBits := float64(math.Ilogb(float64(w)))
	wtnLen1 = uint32(math.Ceil(float64(hashSize) * 8 / wBits))
	wtnLen2 = uint32(math.Floor(math.Log2(float64(wtnLen1)*(w-1))/wBits)) + 1
	wtnLen = wtnLen1 + wtnLen2
}

// returns the hash function to use
func HashFunc() hash.Hash {
	return sha3.New256()
}
