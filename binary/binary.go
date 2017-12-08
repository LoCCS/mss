package binary

func PutUint(b []byte, v uint64) {
	for i := len(b) - 1; i >= 0; i-- {
		b[i] = byte(v & 0xff)
		v >>= 8
	}
}
