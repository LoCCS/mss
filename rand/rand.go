package rand

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/sammy00/mss/config"
	"golang.org/x/crypto/sha3"
)

// Rand produces a random randOTS each time by
//	updating
//		randOTS=hash(seed)
//		seed=hash(seed||randOTS)
type Rand struct {
	seed []byte
}

func New(seed []byte) *Rand {
	rng := new(Rand)
	rng.Seed(seed)
	return rng
}

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

func (rng *Rand) Seed(newSeed []byte) {
	rng.seed = make([]byte, config.Size)
	copy(rng.seed, newSeed)
}

func (rng *Rand) String() string {
	return fmt.Sprintf("{ seed: %s }", hex.EncodeToString(rng.seed))
}

func RandSeed() ([]byte, error) {
	seed := make([]byte, config.Size)
	_, err := rand.Read(seed)

	return seed, err
}
