package main

import (
	"encoding/hex"
	"fmt"

	"github.com/sammy00/mss/config"
	"github.com/sammy00/mss/rand"
)

func main() {
	seed, err := rand.RandSeed()
	if nil != err {
		fmt.Println(err)
		return
	}

	rng := rand.New(seed)
	p := make([]byte, config.Size)
	for i := 0; i < 3; i++ {
		rng.Read(p)

		fmt.Println("rng=", rng)
		fmt.Println("p=", hex.EncodeToString(p))
	}
}
