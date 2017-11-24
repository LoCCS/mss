package mss

import "github.com/sammy00/mss/config"

// merge estimates the hash for (hashLef||hashRight)
func merge(hashLeft, hashRight []byte) []byte {
	hashFunc := config.HashFunc()

	hashFunc.Reset()
	hashFunc.Write(hashLeft)
	hashFunc.Write(hashRight)

	return hashFunc.Sum(nil)
}
