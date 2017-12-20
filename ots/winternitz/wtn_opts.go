package winternitz

import (
	"crypto/rand"
)

// WtnOpts provides a container specifying options for W-OTS
type WtnOpts struct {
	addr  address
	nonce []byte
}

// NewWtnOpts makes a WtnOpts of the specified security level
// and securityLevel should be in {32, 64}
func NewWtnOpts(securityLevel uint32) *WtnOpts {
	nonce := make([]byte, securityLevel)
	if _, err := rand.Read(nonce); nil != err {
		return nil
	}

	return &WtnOpts{newAddress(), nonce}
}

// Clone makes a copy of this WtnOpts
func (opts *WtnOpts) Clone() *WtnOpts {
	optsC := new(WtnOpts)

	// copy of address
	if nil != opts.addr {
		addr := make([]byte, opts.addr.Len())
		copy(addr, []byte(opts.addr))
		optsC.addr = addr
	}

	// copy of nonce
	if nil != opts.nonce {
		optsC.nonce = make([]byte, len(opts.nonce))
		copy(optsC.nonce, opts.nonce)
	}

	return optsC
}

// SetKeyIdx sets the index of the key-pair in use
func (opts *WtnOpts) SetKeyIdx(i uint32) {
	opts.addr.setKeyIdx(i)
}

// SetNonce sets the nonce for this WtnOpts
func (opts *WtnOpts) SetNonce(nonce []byte) {
	copy(opts.nonce, nonce)
}

// Nonce returns the nonce employed by this WtnOpts
func (opts *WtnOpts) Nonce() []byte {
	nonce := make([]byte, len(opts.nonce))
	copy(nonce, opts.nonce)

	return nonce
}

// SecurityLevel returns the security level specified by
// this WtnOpts, should be the same length as
// the nonce in use
func (opts *WtnOpts) SecurityLevel() uint32 {
	return uint32(len(opts.nonce))
}

//Serialize encodes the winternitz options
func (opts *WtnOpts) Serialize() []byte {
	addrLen := uint8(len(opts.addr))
	nonceLen := uint8(len(opts.nonce))

	optsBytes := make([]byte, 1 + addrLen + 1 + nonceLen)
	optsBytes[0] = addrLen
	copy(optsBytes[1:], opts.addr)
	offset := 1 + addrLen
	optsBytes[offset] = nonceLen
	copy(optsBytes[offset + 1:], opts.nonce)
	return optsBytes
}

//Deserialize recovers the WtnOpts from bytes
func Deserialize(optsBytes []byte) *WtnOpts{
	var addr address
	addrLen := optsBytes[0]
	addr = optsBytes[1: 1+addrLen]
	offset := 1 + addrLen
	nonceLen := optsBytes[offset]
	nonce := optsBytes[offset + 1: offset + 1 + nonceLen]
	opts := &WtnOpts{
		addr:addr,
		nonce:nonce,
	}
	return opts
}
// DummyWtnOpts is a dummy WtnOpts with a random nonce
var DummyWtnOpts *WtnOpts = NewWtnOpts(SecurityLevel)
