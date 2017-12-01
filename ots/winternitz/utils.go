package winternitz

import (
	"bytes"
)

// IsEqual checks if two secret key are equal
func IsEqual(sk1, sk2 *PrivateKey) bool {
	for i := range sk1.x {
		if !bytes.Equal(sk1.x[i], sk2.x[i]) {
			return false
		}
	}

	return true
}

// HashPk computes the hash value for a W-OTS public key
func HashPk(pk *PublicKey) []byte {
	h := HashFunc()

	for i := range pk.Y {
		h.Write(pk.Y[i])
	}

	return h.Sum(nil)
}
