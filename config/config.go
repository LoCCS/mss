package config

import (
	"hash"

	"golang.org/x/crypto/sha3"
)

const (
	Size256 = 32
	Size512 = 64
	Size    = Size256
)

func HashFunc() hash.Hash {
	return sha3.New256()
}
