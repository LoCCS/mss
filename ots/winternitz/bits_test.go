package winternitz

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

// TestBitMask tests the format of the mask,
//	which shoulde be in form of `11...1` of length W
func TestBitMask(t *testing.T) {
	mask := bitMask()
	want := []byte{0xff, 0xff}

	if !bytes.Equal(mask.Bytes(), want) {
		t.Fatal("wants %s, got %s", hex.EncodeToString(want),
			hex.EncodeToString(mask.Bytes()))
	}
}

// TestSplit tests the splitting of byte slice
func TestSplit(t *testing.T) {
	digest := []byte{0x00, 0x01, 0x02, 0x03}

	if blocks, err := split(digest, 7); nil == err {
		for _, block := range blocks {
			fmt.Print(block.Text(2), " ")
		}
		fmt.Println()
	} else {
		t.Fatal(err)
	}

	digestInt := new(big.Int)
	digestInt.SetBytes(digest)

	fmt.Println(digestInt.Text(2))
}

// TestHashToBlocks tests encoding of hash digest into blocks for signing
func TestHashToBlocks(t *testing.T) {
	digest := []byte{0x00, 0x01, 0x02, 0x03}

	if blocks, err := hashToBlocks(digest); nil == err {
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
