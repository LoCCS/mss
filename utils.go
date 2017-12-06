package mss

import "github.com/LoCCS/mss/config"
import wots "github.com/LoCCS/mss/ots/winternitz"

// merge estimates the hash for (hashLef||hashRight)
func merge(hashLeft, hashRight []byte) []byte {
	hashFunc := config.HashFunc()

	hashFunc.Reset()
	hashFunc.Write(hashLeft)
	hashFunc.Write(hashRight)

	return hashFunc.Sum(nil)
}

// HashPk computes the hash value for a W-OTS public key
func HashPk(pk *wots.PublicKey) []byte {
	h := config.HashFunc()

	for i := range pk.Y {
		h.Write(pk.Y[i])
	}

	return h.Sum(nil)
}
