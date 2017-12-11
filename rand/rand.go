package rand

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/sha3"

	"github.com/LoCCS/mss/config"
)

// Reader is a globally accessible PRNG (pseudo random number generator) instance
var Reader io.Reader

// init initializes some relevant parameters
func init() {
	seed, _ := RandSeed()
	Reader = New(seed)
}

// RandSeed generates a random seed of predefined bytes
func RandSeed() ([]byte, error) {
	seed := make([]byte, config.Size)
	_, err := rand.Read(seed)

	return seed, err
}

// Rand implements a PRNG based on a sha3.ShakeHash
//	every read from which will update the internal seed as
//	seed := hash(seed)
type Rand struct {
	seed []byte // state seed
	//seedOTS []byte // one-time seed
	sha sha3.ShakeHash
}

// New makes PRNG instance based on the given seed
func New(seed []byte) *Rand {
	rng := new(Rand)
	rng.Seed(seed)

	rng.sha = sha3.NewShake256()

	return rng
}

// Read reads out len(p) random bytes
//	and update the underlying state seed
func (rng *Rand) Read(p []byte) (int, error) {
	rng.sha.Reset()
	rng.sha.Write(rng.seed)

	// update seed
	rng.sha.Read(rng.seed)
	// read out random bytes
	return rng.sha.Read(p)
}

//  ExportSeed exports the seed for next generation
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
