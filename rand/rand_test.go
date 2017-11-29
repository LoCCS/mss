package rand

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/sammy00/mss/config"
)

// TestRandSeed tests the generation of random seeds
func TestRandSeed(t *testing.T) {
	for i := 0; i < 10; i++ {
		seed, err := RandSeed()
		if nil != err {
			t.Fatal(err)
		}

		if len(seed) != config.Size {
			t.Errorf("invalid lenght of seed, wants %v, got %v", config.Size, len(seed))
		}
	}
}

// TestRand tests the correctness of Rand
func TestRand(t *testing.T) {
	seed, err := RandSeed()
	if nil != err {
		fmt.Println(err)
		return
	}

	rng := New(seed)
	p := make([]byte, config.Size)
	rng.Read(p)

	rng2 := New(rng.ExportSeed())
	p2 := make([]byte, config.Size)
	for i := 0; i < 2; i++ {
		rng.Read(p)
		rng2.Read(p2)

		if !bytes.Equal(p, p2) {
			t.Fatalf("wants %s, got %s", hex.EncodeToString(p), hex.EncodeToString(p2))
		}
	}
}
