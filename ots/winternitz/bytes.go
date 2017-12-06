package winternitz

import (
	"math"
)

// ToBaseW outputs an array `out` of integers between 0 and (base - 1)
//	len(out) is REQUIRED to be <=8*len(X)/log2(base)
//	and base should be either 4 or 16
func ToBaseW(out []byte, X []byte, base byte) {

	mask := base - 1
	baseBits := uint32(math.Ilogb(float64(base)))

	consumed := len(out) - 1 // index of out byte filled already
	for i := len(X) - 1; (i >= 0) && (consumed >= 0); i-- {
		for offset := uint32(0); (offset < 8) && (consumed >= 0); offset += baseBits {
			out[consumed] = byte((X[i] >> offset) & mask)
			consumed--
		}
	}
}

// ToBytes estimates the len(out)-byte slice containing the binary
//	representation of x in big-endian byte-order
func ToBytes(out []byte, x uint64) {
	for i := len(out) - 1; i >= 0; i-- {
		out[i] = byte(x & 0xff)
		x >>= 8
	}
}

// hashToBlocks encodes a given hash value as wtnLen base-w blocks
func hashToBlocks(hash []byte) []byte {
	blocks := make([]byte, wtnLen)

	// convert hash to base-w blocks
	ToBaseW(blocks[:wtnLen1], hash, w)

	// compute checksum
	var checksum uint64
	for _, b := range blocks {
		checksum += w - 1 - uint64(b)
	}

	// ?? convert checksum to base-w
	// left shift checksum
	checksum <<= (8 - (wtnLen2 * uint32(math.Ilogb(w)) % 8))
	// big-endian-order byte string of checksum
	checksumLen := int(math.Ceil(float64(wtnLen2) * float64(math.Ilogb(float64(w))) / 8))
	checksumBytes := make([]byte, checksumLen)
	ToBytes(checksumBytes, checksum)
	ToBaseW(blocks[wtnLen1:], checksumBytes, w)

	return blocks
}
