package winternitz

// evalChain computes an iteration of fSum() on an n-byte input
//	using outputs of prf()
func evalChain(x []byte, offset, numIter uint32, nonce []byte, addr address) []byte {
	if (offset + numIter) > wtnMask {
		return nil
	}

	out := make([]byte, len(nonce))
	copy(out, x)

	key := make([]byte, len(nonce))
	bitmask := make([]byte, len(nonce))
	for i := offset; i < (offset + numIter); i++ {
		addr.setHashAddress(i)

		// derive key
		addr.onKey()
		prf(key, nonce, addr)
		// derive bitmask
		addr.onMask()
		prf(bitmask, nonce, addr)

		// out ^ bitmask
		for j := range out {
			out[j] ^= bitmask[j]
		}

		// advance to next hash
		fSum(out, key, out)
	}

	return out
}
