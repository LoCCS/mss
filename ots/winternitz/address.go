package winternitz

import "fmt"

// address encodes an OTS hash address, which is formatted as
/*
	+-----------------------+
	| layer address	(32 bit)|
	+-----------------------+
	| tree address	(64 bit)|
	+-----------------------+
	| type = 0			(32 bit)|
	+-----------------------+
	| OTS address		(32 bit)|
	+-----------------------+
	| chain address (32 bit)|
	+-----------------------+
	| hash address	(32 bit)|
	+-----------------------+
	| keyAndMask		(32 bit)|
	+-----------------------+
*/
type address []byte

// byte offset of different component within address
const (
	addr_type_offset         = 4 + 8
	addr_ots_offset          = addr_type_offset + 4
	addr_chain_offset        = addr_ots_offset + 4
	addr_hash_offset         = addr_chain_offset + 4
	addr_key_and_mask_offset = addr_hash_offset + 4
)

// newAddress make a OTS address with all bytes as 0
func newAddress() address {
	raw := make([]byte, 32)

	// layer address as 0: raw[0:4]		= 0x00 00 00 00
	//  tree address as 0: raw[4:12]	= 0x00 00 00 00 00 00 00 00 00
	//		type field as 0: raw[12:16]	= 0x00 00 00 00

	return address(raw)
}

// onKey masks the last 4 bytes as all 0s
//	to prepare the address for generating keys
func (addr address) onKey() {
	bs := []byte(addr)

	ToBytes(bs[addr_key_and_mask_offset:], 0)
}

// onMask masks the last 4 bytes as all 1s
//	to prepare the address for generating bitmasks
func (addr address) onMask() {
	bs := []byte(addr)

	ToBytes(bs[addr_key_and_mask_offset:], 0xffffffff)
}

// setChainAddress sets the index i of target component sk[i]
func (addr address) setChainAddress(chainAddress uint32) {
	bs := []byte(addr)

	ToBytes(bs[addr_chain_offset:addr_hash_offset], uint64(chainAddress))
}

// setHashAddress set the offset of the address
//	w.r.t the hash chain starting from the private key
func (addr address) setHashAddress(hashAddress uint32) {
	bs := []byte(addr)

	ToBytes(bs[addr_hash_offset:addr_key_and_mask_offset], uint64(hashAddress))
}

// setKeyIdx encodes the index of the OTS key pair within the tree
func (addr address) setKeyIdx(i uint32) {
	bs := []byte(addr)

	ToBytes(bs[addr_ots_offset:addr_chain_offset], uint64(i))
}

// Len returns the length in bytes of the address
func (addr address) Len() int {
	return len([]byte(addr))
}

// String returns the encoded hex string of the address
func (addr address) String() string {
	return fmt.Sprintf("%x", []byte(addr))
}
