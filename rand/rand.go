package rand

import (
	"crypto/rand"
	"fmt"

	"github.com/merkle-signature-scheme/config"
	"golang.org/x/crypto/sha3"
)

// RandomSeed returns a random seed
func RandomSeed() ([]byte, error) {
	seed := make([]byte, config.N)
	_, err := rand.Read(seed)

	return seed, err
}

// PRNG returns seedDot and next seed
func PRNG(seed []byte) ([]byte, []byte, error) {
	if len(seed) != config.N {
		e := fmt.Errorf("seed must be a []bytes of %d length", config.N)
		return []byte{}, []byte{}, e
	}
	// seedDot=hash(seed)
	seedDot := sha3.Sum256(seed) // [32]byte

	// nextSeed=hash(seed || seedDot)
	source := append(seed, seedDot[:]...)
	nextSeed := sha3.Sum256(source)

	return seed[:], nextSeed[:], nil
}
