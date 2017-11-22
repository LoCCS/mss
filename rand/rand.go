package rand

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/sammy00/mss/config"
	"golang.org/x/crypto/sha3"
)

// Reader is a globally accessible PRNG instance
var Reader io.Reader

// init initializes some relevant parameters
func init() {
	seed, _ := RandSeed()
	Reader = New(seed)
}

// Rand produces a random randOTS each time by
//	updating
//		randOTS=hash(seed)
//		seed=hash(seed||randOTS)
type Rand struct {
	seed []byte
}

// New makes PRNG instance based on the given seed
func New(seed []byte) *Rand {
	rng := new(Rand)
	rng.Seed(seed)
	return rng
}

// Read reads min(len(p),config.Size) random bytes
//	from the PRNG
func (rng *Rand) Read(p []byte) (int, error) {
	// update randOTS
	randOTS := sha3.Sum256(rng.seed)
	copy(p, randOTS[:])

	// update seed
	randOTS = sha3.Sum256(append(rng.seed, randOTS[:]...))
	copy(rng.seed, randOTS[:])

	n := len(p)
	if n > len(randOTS) {
		n = len(randOTS)
	}

	return n, nil
}

// Seed reset the PRNG to be seeded at the new seed
func (rng *Rand) Seed(newSeed []byte) {
	rng.seed = make([]byte, config.Size)
	copy(rng.seed, newSeed)
}

// String outputs the string representation of this PRNG
func (rng *Rand) String() string {
	return fmt.Sprintf("{ seed: %s }", hex.EncodeToString(rng.seed))
}

// RandSeed generates a random seed
func RandSeed() ([]byte, error) {
	seed := make([]byte, config.Size)
	_, err := rand.Read(seed)

	return seed, err
}
