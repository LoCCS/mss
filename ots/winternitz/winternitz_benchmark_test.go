package winternitz

import (
	"testing"

	mrand "github.com/LoCCS/mss/rand"
)

func BenchmarkGenerateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateKey(mrand.Reader)
	}
}
