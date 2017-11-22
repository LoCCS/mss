package winternitz

import (
	"errors"
	"math"
	"math/big"

	"github.com/sammy00/mss/config"
)

// predefined parameters determined by
//	* bit length `config.Size` of the hash function in use
//	* the number bits `W` to manipulate simultaneously
var (
	t1 = int(math.Ceil(float64(config.Size*8) / float64(W)))
	t2 = int(math.Ceil(float64(math.Ilogb(float64(t1))+1+W) / float64(W)))
	t  = t1 + t2
)

// bitMask makes a bit mask as `11...1` of length W
func bitMask() *big.Int {
	mask := big.NewInt(1)
	return mask.Lsh(mask, W).Sub(mask, big.NewInt(1))
}

// pow2ToW returns 2^W
func pow2ToW() *big.Int {
	twoToW := big.NewInt(1)
	return twoToW.Lsh(twoToW, W)
}

// split partitions a given digest into t blocks {b_i},
//  each of which is of W bits
func split(digest []byte, t int) ([]*big.Int, error) {
	// convert digest as a big-endian byte slice into a big integer
	bs := new(big.Int)
	bs.SetBytes(digest)

	//fmt.Println("t=", t, "W=", W, "bs.BitLen()=", bs.BitLen())
	if t*W < bs.BitLen() {
		//fmt.Println(t*W, len(digest)*8)
		return nil, errors.New("invalid number of blocks")
	}

	// split bs into t blocks
	mask := bitMask()
	blocks := make([]*big.Int, t)
	for i := len(blocks) - 1; i >= 0; i-- {
		blocks[i] = big.NewInt(0)
		//fmt.Println(blocks[i])
		blocks[i].And(mask, bs)
		bs.Rsh(bs, W)
	}

	return blocks, nil
}

// checksum estimates the checksum on {b_i} as
//	c = \sum_{i=t-t_1}^{t-1}(2^W-b_i)
func checksum(blocks []*big.Int) *big.Int {
	// 2^W
	twoToW := pow2ToW()

	c := big.NewInt(0)
	for _, block := range blocks {
		c.Add(c, twoToW).Sub(c, block)
	}

	return c
}

// hashToBlocks encodes a given hash value into  t blocks
func hashToBlocks(hash []byte) ([]*big.Int, error) {
	// blocks split from digest
	blocks, err := split(hash, t1)
	if nil != err {
		return nil, err
	}

	//fmt.Println("**after 1st split")
	// blocks split from checksum
	blocksC, err := split(checksum(blocks).Bytes(), t2)
	if nil != err {
		return nil, err
	}

	return append(blocks, blocksC...), nil
}
