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

// ToBytes estimates the y-byte slice containing the binary
//	representation of X in big-endian byte-order
func ToBytes(X uint64) []byte {
	out := make([]byte, 8)

	// bits to shift
	var offset uint32
	//fmt.Println("len(out)-1:", len(out)-1)
	for i := len(out) - 1; i >= 0; i-- {
		out[i] = byte((X >> offset) & 0xff)
		offset += 8
	}

	return out
}

// hashToBlocks encodes a given hash value as t base-w blocks
func hashToBlocks(hash []byte) []byte {
	blocks := make([]byte, t)

	//fmt.Println("t =", t)
	//fmt.Println("len1 =", len1)
	//fmt.Println("len2 =", len2)

	// convert hash to base-w blocks
	ToBaseW(blocks[:len1], hash, w)

	// compute checksum
	var checksum uint64
	for _, b := range blocks {
		//fmt.Println(uint64(b))
		checksum += w - 1 - uint64(b)
	}

	// ?? convert checksum to base-w
	//checksum <<= (8 - (len2 * uint32(math.Ilogb(w)) % 8))
	checksumBytes := ToBytes(checksum)
	ToBaseW(blocks[len1:], checksumBytes, w)

	return blocks
}
