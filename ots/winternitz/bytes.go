package winternitz

import "math"

// ToBaseW outputs an array `out` of integers between 0 and (base - 1)
//	len(out) is REQUIRED to be <=8*len(X)/log2(base)
//	and base should be either 4 or 16
func ToBaseW(out []byte, X []byte, base byte) {

	mask := base - 1
	baseBits := uint32(math.Ilogb(float64(base)))

	consumed := len(out) - 1 // index of out byte filled already
	for i := len(X) - 1; i >= 0; i-- {
		for offset := uint32(0); offset < 8; offset += baseBits {
			out[consumed] = byte((X[i] >> offset) & mask)
			consumed--
		}
	}
}
