package rand

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/sammy00/mss/config"
)

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

func TestRand(t *testing.T) {
	seed, err := RandSeed()
	if nil != err {
		fmt.Println(err)
		return
	}

	rng := New(seed)
	p := make([]byte, config.Size)
	for i := 0; i < 10; i++ {
		rng.Read(p)
		fmt.Println("p=", hex.EncodeToString(p))
	}
}
