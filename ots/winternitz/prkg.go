package winternitz

import (
	"github.com/LoCCS/mss/rand"
)

// KeyIterator is a iterator to produce a key chain for
//	user based on a seed
type KeyIterator struct {
	rng *rand.Rand
	// the 0-based index of next running iteration
	//	w.r.t the initial genesis seed
	offset uint32
	// options specifying stuff like nonce for
	//	randomizing hash function
	//*WtnOpts
}

// NewKeyIterator makes a key pair iterator
func NewKeyIterator(compactSeed []byte) *KeyIterator {
	return &KeyIterator{
		rng:    rand.New(compactSeed),
		offset: 0,
	}
}

// Init resets the KeyIterator
func (iter *KeyIterator) Init(compositeSeed []byte) bool {
	seedLen := len(compositeSeed) - 4
	if seedLen <= 0 {
		return false
	}

	iter.rng = rand.New(compositeSeed[:seedLen])
	iter.offset = uint32(GetUint64(compositeSeed[seedLen:]))

	return true
}

// Next estimates and returns the next sk-pk pair
func (iter *KeyIterator) Next() (*PrivateKey, error) {
	iter.offset++
	keyPair, err := GenerateKey(DummyWtnOpts, iter.rng)
	iter.rng.NextState()
	return keyPair, err
}

// Offset returns 0-based index of the **next** running iteration
func (iter *KeyIterator) Offset() uint32 {
	return iter.offset
}

// Seed returns the internal updated seed for usage
//	such as saving state of the iterator
// !!!TBR: to be removed
func (iter *KeyIterator) Seed() []byte {
	return iter.rng.ExportSeed()
}

// Serialize encodes the key iterator as a integrated seed
//	in form of seed||offset
func (iter *KeyIterator) Serialize() []byte {
	seed := iter.rng.ExportSeed()
	seedLen := len(seed)
	// append 4 bytes to the end to make space for the offset
	for i := 0; i < 4; i++ {
		seed = append(seed, byte(0))
	}

	//binary.BigEndian.PutUint32(seed[seedLen:], iter.offset)
	ToBytes(seed[seedLen:], uint64(iter.offset))

	return seed
}
