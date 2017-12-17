package winternitz

import "encoding/binary"

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
