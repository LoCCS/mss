package winternitz

import (
	"errors"

	"golang.org/x/crypto/sha3"
)

// extra pre-padding in big-endian order to
//	+ adapt input for compress function
//	+ achieve independence among different function families
var paddings [][]byte

func init() {
	paddings = make([][]byte, 4)

	for i := range paddings {
		paddings[i] = make([]byte, SecurityLevel)
		paddings[i][SecurityLevel-1] = byte(i)
	}
}

// keyedHashSum implements a keyed hash functions
//			 taking input: padding(64B)||key||data
//	generating output to buffer hash
//	report error if any errors
func keyedHashSum(hash, key, data []byte, paddingType int) error {
	if paddingType >= len(paddings) {
		return errors.New("invalid paddingType")
	}

	sh := sha3.NewShake256()

	// write input as paddings[paddingType]||key||data
	sh.Write(paddings[paddingType])
	sh.Write(key)
	sh.Write(data)

	// read out output
	sh.Read(hash)

	return nil
}

// fSum estimates keyed hash for input (padding00||key||data)
func fSum(hash, key, data []byte) {
	keyedHashSum(hash, key, data, 0x00)
}

// prf estimates keyed hash for input (padding03||key||data)
func prf(hash, key, data []byte) {
	keyedHashSum(hash, key, data, 0x03)
}
