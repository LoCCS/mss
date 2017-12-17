package winternitz

import (
	"encoding/binary"
)

// ToBase outputs an array `out` of integers between 0 and ((2<<baseWidth) - 1)
// len(out) is REQUIRED to be <=8*len(X)/baseWidth
// and baseWidth should be a member in set {1,2,4,8}.
// Otherwise, the result is unpredictable
func ToBase(out []byte, X []byte, baseWidth uint8) {
	mask := byte((1 << baseWidth) - 1)

	// length of output
	ell := len(out)
	consumed := 0
	for _, x := range X {
		// a do-while loop
		for bits := 8 - baseWidth; ; bits -= baseWidth {
			if consumed >= ell {
				return // no more buffer
			}

			out[consumed] = (x >> bits) & mask
			consumed++

			if 0 == bits {
				break
			}
		}
	}
}

// hashToBlocks encodes a given hash value as wtnLen base-w blocks
func hashToBlocks(hash []byte) []byte {
	blocks := make([]byte, wtnLen)

	// convert hash to base-w blocks
	ToBase(blocks[:wtnLen1], hash, w)

	// compute checksum
	var checksum uint16
	for _, b := range blocks {
		// + (2^w-1)+b[i]
		checksum += wtnMask - uint16(b)
	}

	// left shift checksum
	checksum <<= (8 - (wtnLen2*w)%8)
	// big-endian-order byte string of checksum
	var cb [2]byte
	binary.BigEndian.PutUint16(cb[:], checksum)
	ToBase(blocks[wtnLen1:], cb[:], w)

	return blocks
}
