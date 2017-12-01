package winternitz

// HashChainEval computes an iteration of the hash function
//	defined by HashFunc() on an input byte slice
func HashChainEval(in []byte, numIter uint32) []byte {

	out := make([]byte, len(in))
	copy(out, in)

	h := HashFunc()
	for i := numIter; i > 0; i-- {
		// out = hash(out)
		h.Reset()
		h.Write(out)
		out = h.Sum(nil)
	}

	return out
}
