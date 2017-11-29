package rand

import (
	"crypto/rand"
	"hash"
	"io"

	"github.com/sammy00/mss/config"
	"golang.org/x/crypto/sha3"
)

// sha is the internal hash function to build the PRNG
var sha hash.Hash

// Reader is a globally accessible PRNG (pseudo random number generator) instance
var Reader io.Reader

// init initializes some relevant parameters
func init() {
	// initialize sha as sha3
	sha = sha3.New256()

	seed, _ := RandSeed()
	Reader = New(seed)
}

// RandSeed generates a random seed of predefined bytes
func RandSeed() ([]byte, error) {
	seed := make([]byte, sha.Size())
	_, err := rand.Read(seed)

	return seed, err
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

// Read reads min(len(p),config.Size) random bytes from the PRNG
func (rng *Rand) Read(p []byte) (int, error) {
	sha.Reset()

	// update randOTS
	sha.Write(rng.seed)
	randOTS := sha.Sum(nil)
	sz := copy(p, randOTS)

	// update seed
	sha.Write(randOTS)
	randOTS = sha.Sum(nil)
	copy(rng.seed, randOTS)

	return sz, nil
}

//  ExportSeed exports the updated seed for next generation
func (rng *Rand) ExportSeed() []byte {
	seed := make([]byte, len(rng.seed))
	copy(seed, rng.seed)

	return seed
}

// Seed uses provided seed to initialize the generator to a deterministic state
func (rng *Rand) Seed(seed []byte) {
	if nil == rng.seed {
		rng.seed = make([]byte, config.Size)
	}
	copy(rng.seed, seed)
}
