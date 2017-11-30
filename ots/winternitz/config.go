package winternitz

import (
	"math"

	"github.com/sammy00/mss/config"
)

const (
	w = 16 // the Winternitz parameter, it is a member of the set {4, 16}
)

// length in bytes of digest produced by the employed hash function
var hashSize int

// parameters derived from w
// skLen is the length of secret key in bytes
var len1, len2, skLen uint32

func init() {
	hashSize = config.HashFunc().Size()

	wBits := float64(math.Ilogb(float64(w)))
	len1 = uint32(math.Ceil(float64(hashSize) * 8 / wBits))
	//len2 = uint32(math.)
}
