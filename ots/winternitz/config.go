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
	w = 4 // the width in bits of the Winternitz parameter, it should be in {2, 4}
)

// parameters derived from w
//	wtnLen is the number of n-byte string elements
//	in a WOTS+ private key, public key, and signature.
var wtnLen1, wtnLen2, wtnLen uint32

func init() {
	hashSize = HashFunc().Size()

	wtnLen1 = uint32(SecurityLevel * 8 / w)
	wtnLen2 = uint32(math.Floor(math.Log2(float64(wtnLen1*((1<<w)-1)))/w)) + 1

	wtnLen = wtnLen1 + wtnLen2
	//fmt.Println("****", wtnLen1, wtnLen2, wtnLen)
}

// returns the hash function to use
func HashFunc() hash.Hash {
	return sha3.New256()
}
