package rand

import (
	"crypto/rand"
	"hash"
	"io"

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

// Rand produces a random seedOTS each time by
//	updating
//		seedOTS_0=hash(seed_i)
//		seed_{i+1}=hash(seed_i||seedOTS_0)
//		seedOTS_{j+1} = hash(seedOTS_j)
type Rand struct {
	seed    []byte // state seed
	seedOTS []byte // one-time seed
}

// New makes PRNG instance based on the given seed
func New(seed []byte) *Rand {
	rng := new(Rand)
	rng.Seed(seed)

	return rng
}

// NextState updates the state seed after updating, we should have
//		* seedOTS_0	= hash(seed_i)
//		* seed_{i+1}= hash(seed_i||seedOTS_0)
func (rng *Rand) NextState() {
	// seed = hash(seed||hash(seed))
	sha.Reset()
	sha.Write(rng.seed)
	sha.Write(sha.Sum(nil))
	rng.seed = sha.Sum(nil)

	// seedOTS = hash(seed)
	sha.Reset()
	sha.Write(rng.seed)
	rng.seedOTS = sha.Sum(nil)

	// don't forget to clear up the internal state of sha
	sha.Reset()
}

// Read reads min(len(p),config.Size) random bytes from the PRNG
func (rng *Rand) Read(p []byte) (int, error) {
	// next seed for one-time use
	//	seedOTS = hash(seedOTS)
	sha.Reset()
	sha.Write(rng.seedOTS)
	rng.seedOTS = sha.Sum(nil)

	// fill up p
	sz := copy(p, rng.seedOTS)

	// clear up the internal state of sha
	sha.Reset()

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
		rng.seed = make([]byte, sha.Size())
	}

	copy(rng.seed, seed)

	// initialize seedOTS = hash(seed)
	sha.Reset()
	sha.Write(rng.seed)
	rng.seedOTS = sha.Sum(nil)

	// clear the internal state of sha to avoid state leakage
	sha.Reset()
}
