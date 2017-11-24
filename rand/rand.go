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
	hashFunc := sha3.New256()

	hashFunc.Write(rng.seed)
	// update randOTS
	randOTS := hashFunc.Sum(nil)
	sz := copy(p, randOTS)

	// update seed
	hashFunc.Write(randOTS)
	randOTS = hashFunc.Sum(nil)
	copy(rng.seed, randOTS)

	return sz, nil
}

func (rng *Rand) TellMeSeed() []byte {
	seed := make([]byte, len(rng.seed))
	copy(seed, rng.seed)

	return seed
}

// Seed reset the PRNG to be seeded at the new seed
func (rng *Rand) Seed(newSeed []byte) {
	if nil == rng.seed {
		rng.seed = make([]byte, config.Size)
	}
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
