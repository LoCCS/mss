package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"golang.org/x/crypto/sha3"
)

func SHA3Demo() {
	fmt.Println("***demo of SHA3***")
	sha := sha3.New256()
	msg := []byte{0x01, 0x02, 0x03, 0x04}
	digest := sha.Sum(msg)
	fmt.Println(hex.EncodeToString(digest))
	fmt.Println("***end demo of SHA3***")
}

func bigIntDemo() {
	fmt.Println("***demo of big.Int***")
	one := big.NewInt(1)

	fmt.Println(one)
	fmt.Println(one.Sign())

	buf := []byte{0x00, 0x01, 0x02, 0x03, 0x04}
	z := new(big.Int)
	z.SetBytes(buf)

	fmt.Println("z=", hex.EncodeToString(z.Bytes()))
	fmt.Println("***end demo of big.Int***")
}

func main() {
	//SHA3Demo()
	//bigIntDemo()

	hello := big.NewInt(1)
	fmt.Println(hello.Not(hello))
}
