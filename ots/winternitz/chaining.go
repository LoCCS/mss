package winternitz

// evalChain computes an iteration of fSum() on an n-byte input
//	using outputs of prf()
func evalChain(x []byte, offset, numIter uint32, wtnOpts *WtnOpts) []byte {
	if (offset + numIter) >= w {
		return nil
	}

	out := make([]byte, wtnOpts.SecurityLevel())
	copy(out, x)

	key := make([]byte, wtnOpts.SecurityLevel())
	bitmask := make([]byte, SecurityLevel)
	for i := offset; i < (offset + numIter); i++ {
		wtnOpts.addr.setHashAddress(i)

		// derive key
		wtnOpts.addr.onKey()
		prf(key, wtnOpts.nonce, wtnOpts.addr)
		// derive bitmask
		wtnOpts.addr.onMask()
		prf(bitmask, wtnOpts.nonce, wtnOpts.addr)

		// out ^ bitmask
		for j := range out {
			out[j] ^= bitmask[j]
		}

		// advance to next hash
		fSum(out, key, out)
	}

	return out
}
