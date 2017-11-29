package winternitz

import (
	"encoding/json"

	"github.com/sammy00/mss/rand"
)

// KeyIterator is a iterator to produce a key chain for
//	user based on a seed
type KeyIterator struct {
	rng *rand.Rand
	// the 0-based index of next running iteration
	//	w.r.t the initial genesis seed
	offset uint32
}

// NewKeyIterator makes a key pair iterator
func NewKeyIterator(seed []byte) *KeyIterator {
	return &KeyIterator{rand.New(seed), 0}
}

// Init resets the KeyIterator
func (iter *KeyIterator) Init(seed []byte, offset uint32) {
	iter.rng = rand.New(seed)
	iter.offset = offset
}

// Next estimates and returns the next sk-pk pair
func (iter *KeyIterator) Next() (*PrivateKey, error) {
	iter.offset++
	return GenerateKey(iter.rng)
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

// keyIteratorEx is a version of KeyIterator to export
type keyIteratorEx struct {
	Seed   []byte
	Offset uint32
}

// MarshalAsJSON encodes a key iterator in JSON format
func (iter *KeyIterator) MarshalAsJSON() ([]byte, error) {
	return json.Marshal(&keyIteratorEx{
		Seed:   iter.rng.ExportSeed(),
		Offset: iter.offset,
	})
}

// Unmarshal decodes a key iterator from a source byte slice
func (iter *KeyIterator) UnmarshalFromJSON(src []byte) error {
	itrEx := new(keyIteratorEx)
	err := json.Unmarshal(src, itrEx)

	if nil == err {
		iter.Init(itrEx.Seed, itrEx.Offset)
	}

	return err
}
