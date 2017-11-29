package mss_test

import (
	"encoding/json"
	"fmt"
	"testing"
)

type hello struct {
	Who string
}

type world struct {
	Rank   int32
	Prefix hello
}

func TestMSS(t *testing.T) {
	wd := world{
		Rank:   123,
		Prefix: hello{"world"},
	}

	data, err := json.Marshal(wd)
	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println(string(data))
}
