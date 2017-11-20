package winternitz

import (
	"fmt"
	"math/big"
	"testing"
)

func TestSplit(t *testing.T) {
	digest := []byte{0x00, 0x01, 0x02, 0x03}

	if blocks, err := Split(digest, 7); nil == err {
		for _, block := range blocks {
			fmt.Print(block.Text(2), " ")
		}
	} else {
		fmt.Println(err)
	}
	fmt.Println()

	digestInt := big.NewInt(0)
	digestInt.SetBytes(digest)
	fmt.Println(digestInt.Text(2))
}
