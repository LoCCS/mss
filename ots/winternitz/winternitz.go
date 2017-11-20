package winternitz

import (
	"errors"
	"math/big"
)

const (
	W = 18
)

func split(digest []byte, t int) ([]*big.Int, error) {
	if t*W < len(digest)*8 {
		return nil, errors.New("invalid number of blocks")
	}

	// convert digest as a big-endian byte slice into a big integer
	digestInt := new(big.Int)
	digestInt.SetBytes(digest)

	// split digestInt into t blocks

	return nil, nil
}
