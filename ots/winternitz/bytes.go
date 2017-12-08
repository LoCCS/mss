package winternitz

import (
	"math"
)

// GetUint64 decodes a uint64 from buf in big-endian order
//	if len(buf)>4, only the first 4 bytes will be decoded
func GetUint64(buf []byte) uint64 {
	var x uint64

	for i := range buf {
		if i >= 8 {
			break
		}

		x = (x << 8) | uint64(buf[i])
	}

	return x
}

// ToBase outputs an array `out` of integers between 0 and ((2<<baseWidth) - 1)
//	len(out) is REQUIRED to be <=8*len(X)/baseWidth
//	and baseWidth is a member in set {1,2,4,8}
func ToBase(out []byte, X []byte, baseWidth uint8) {
	mask := byte((1 << baseWidth) - 1)

	consumed := len(out) - 1 // the smallest index of out byte filled already
	for i := len(X) - 1; (i >= 0) && (consumed >= 0); i-- {
		for offset := uint8(0); (offset < 8) && (consumed >= 0); offset += baseWidth {
			out[consumed] = (X[i] >> offset) & mask
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
	ToBase(blocks[:wtnLen1], hash, w)

	// compute checksum
	var checksum uint64
	for _, b := range blocks {
		// + (2^w-1)+b[i]
		checksum += wtnMask - uint64(b)
	}

	// ?? convert checksum to base-w
	// left shift checksum
	checksum <<= (8 - (wtnLen2*w)%8)
	// big-endian-order byte string of checksum
	checksumLen := int(math.Ceil(float64(wtnLen2*w) / 8))
	checksumBytes := make([]byte, checksumLen)
	ToBytes(checksumBytes, checksum)
	ToBase(blocks[wtnLen1:], checksumBytes, w)

	return blocks
}
