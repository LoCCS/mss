package winternitz

import (
	"math"
)

// security parameter
const (
	SecurityLevel256 = 32 // 256 bits
	SecurityLevel512 = 64 // 512 bits
	SecurityLevel    = SecurityLevel256
)

const (
	w       = 4 // the width in bits of the Winternitz parameter, it should be in {2, 4}
	wtnMask = (1 << w) - 1
)

// parameters derived from w
//	wtnLen is the number of n-byte string elements
//	in a WOTS+ private key, public key, and signature.
var wtnLen1, wtnLen2, wtnLen uint32

func init() {
	// wtnLen1 = ceil(8n/w)
	wtnLen1 = uint32(SecurityLevel * 8 / w)

	// wtnLen2 = floor(log2(wtnLen1*(2^w-1))/w)+1
	wtnLen2 = uint32(math.Floor(math.Log2(float64(wtnLen1*wtnMask))/w)) + 1

	wtnLen = wtnLen1 + wtnLen2
}
