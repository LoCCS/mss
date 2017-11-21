package rand

import (
	"reflect"
	"testing"

	"github.com/sammy00/mss/config"
)

func TestRandomSeed(t *testing.T) {
	r1, _ := RandomSeed()
	r2, _ := RandomSeed()
	r3, _ := RandomSeed()
	if len(r1) != config.Size {
		t.Error("length or random seed != %d", config.Size)
	}
	if reflect.DeepEqual(r1, r2) || reflect.DeepEqual(r2, r3) || reflect.DeepEqual(r3, r2) {
		t.Error("Random seeds are equal")
	}
}

func TestPRNG(t *testing.T) {
	r1, _ := RandomSeed()
	seedDot1, nextSeed1, _ := PRNG(r1[:])

}
