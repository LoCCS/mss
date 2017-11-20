package winternitz

import (
	"errors"
	"math/big"
)

const (
	W = 5
)

var oneMask *big.Int

func init() {
	oneMask = big.NewInt(1)
	one := big.NewInt(1)
	oneMask.Lsh(one, W).Sub(oneMask, one)

	//fmt.Println("one=", oneMask.Text(16))
}

// split partitions a given digest into t blocks,
// each of which is of W bits
func split(digest []byte, t int) ([]*big.Int, error) {
	if t*W < len(digest)*8 {
		return nil, errors.New("invalid number of blocks")
	}

	// convert digest as a big-endian byte slice into a big integer
	digestInt := new(big.Int)
	digestInt.SetBytes(digest)

	// split digestInt into t blocks
	blocks := make([]*big.Int, t)
	for i := len(blocks) - 1; i >= 0; i-- {
		blocks[i] = big.NewInt(0)
		//fmt.Println(blocks[i])
		blocks[i].And(oneMask, digestInt)
		digestInt.Rsh(digestInt, W)
	}

	return blocks, nil
}
