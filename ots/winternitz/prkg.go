package winternitz

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/LoCCS/mss/rand"
)

// KeyIterator is a prkgator to produce a key chain for
//	user based on a seed
type KeyIterator struct {
	rng *rand.Rand
	// the 0-based index of next running prkgation
	//	w.r.t the initial genesis seed
	offset uint32
	// options specifying stuff like nonce for
	//	randomizing hash function
	*WtnOpts
}

// NewKeyIterator makes a key pair prkgator
func NewKeyIterator(compactSeed []byte) *KeyIterator {
	prkg := new(KeyIterator)

	prkg.rng = rand.New(compactSeed)
	prkg.offset = 0
	prkg.WtnOpts = NewWtnOpts(SecurityLevel)

	return prkg
}

// Init initialises the prkgator with the composite seed
//	exported by Serialize()
func (prkg *KeyIterator) Init(compositeSeed []byte) bool {
	buf := bytes.NewBuffer(compositeSeed)

	var fieldLen uint8
	// 1. len(seed)
	if err := binary.Read(buf, binary.BigEndian,
		&fieldLen); (nil != err) && (0 == fieldLen) {
		return false
	}
	// 2. compactSeed
	compactSeed := make([]byte, fieldLen)
	if err := binary.Read(buf, binary.BigEndian,
		compactSeed); nil != err {
		return false
	}
	// initialise rng
	prkg.rng = rand.New(compactSeed)

	// initialise WtnOpts if needed before going on
	if nil == prkg.WtnOpts {
		prkg.WtnOpts = NewWtnOpts(SecurityLevel)
	}

	// 3. offset
	var offset uint32
	if err := binary.Read(buf, binary.BigEndian,
		&offset); nil != err {
		return false
	}
	// feed offset to WtnOpts
	prkg.offset = offset
	prkg.SetKeyIdx(prkg.offset)

	// 4. len(nonce)
	fieldLen = 0
	if err := binary.Read(buf, binary.BigEndian,
		&fieldLen); (nil != err) && (0 == fieldLen) {
		return false
	}
	// 5. nonce
	nonce := make([]byte, fieldLen)
	if err := binary.Read(buf, binary.BigEndian,
		nonce); io.ErrUnexpectedEOF == err {
		return false
	}
	// feed nonce to WtnOpts
	prkg.WtnOpts.SetNonce(nonce)

	return true
}

// Next estimates and returns the next sk-pk pair
func (prkg *KeyIterator) Next() (*PrivateKey, error) {
	prkg.WtnOpts.SetKeyIdx(prkg.offset)
	keyPair, err := GenerateKey(prkg.WtnOpts, prkg.rng)

	prkg.offset++

	return keyPair, err
}

// Offset returns 0-based index of the **next** running prkgation
func (prkg *KeyIterator) Offset() uint32 {
	return prkg.offset
}

// Serialize encodes the key iterator as
//	+---------------------------------------------+
//	|	len(seed)||seed||offset||len(nonce)||nonce	|
//	+---------------------------------------------+
//	the byte slice export from here makes up
//	everything needed to recovered the state the prkg
//	So unless it's your first-time use, you should
//	store this byte slice so as to snapshot the prkg
func (prkg *KeyIterator) Serialize() []byte {
	buf := new(bytes.Buffer)

	seed := prkg.rng.ExportSeed()
	// len(seed)
	binary.Write(buf, binary.BigEndian, uint8(len(seed)))
	// seed
	binary.Write(buf, binary.BigEndian, seed)

	// offset
	binary.Write(buf, binary.BigEndian, prkg.offset)

	// len(nonce)
	binary.Write(buf, binary.BigEndian, uint8(prkg.SecurityLevel()))
	// nonce
	binary.Write(buf, binary.BigEndian, prkg.Nonce())

	return buf.Bytes()
}
