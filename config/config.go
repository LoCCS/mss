package config

import (
	"hash"

	"golang.org/x/crypto/sha3"
)

const (
	Size224 = 28
	Size256 = 32
	Size    = Size256
)

func HashFunc() hash.Hash {
	return sha3.New256()
}
