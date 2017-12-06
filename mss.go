package mss

import (
	"bytes"
	"errors"
	"math"

	"github.com/LoCCS/mss/config"
	wots "github.com/LoCCS/mss/ots/winternitz"
	"github.com/LoCCS/mss/rand"
)

// MerkleAgent implements a agent working
//	according to the Merkle signature scheme
type MerkleAgent struct {
	H uint32
	//NumLeafUsed    uint32
	auth           [][]byte
	root           []byte
	seed0          []byte
	keyItr         *wots.KeyIterator
	treeHashStacks []*TreeHashStack
	stackIterators []*wots.KeyIterator
}

// NewMerkleAgent makes a fresh Merkle signing routine
//	by running the generate key and setup procedure
func NewMerkleAgent(H uint32, seed []byte) (*MerkleAgent, error) {
	if H < 2 {
		return nil, errors.New("H should be larger than 1")
	}

	agent := new(MerkleAgent)
	agent.H = H
	agent.auth = make([][]byte, H)
	agent.keyItr = wots.NewKeyIterator(seed)
	agent.seed0 = seed
	agent.treeHashStacks = make([]*TreeHashStack, H)
	agent.stackIterators = make([]*wots.KeyIterator, H)

	// back up the initial seed for later recovery
	integatedSeed := agent.keyItr.Serialize()

	globalStack := NewTreeHashStack(0, H+1)
	for h := uint32(0); h < H; h++ {
		globalStack.Update(1, agent.keyItr)

		agent.treeHashStacks[h] = NewTreeHashStack(0, h)

		agent.treeHashStacks[h].nodeStack.Push(globalStack.Top())
		agent.treeHashStacks[h].SetLeaf(1 << h)
		agent.stackIterators[h] = nil

		globalStack.Update((1<<(h+1))-1, agent.keyItr)
		agent.auth[h] = make([]byte, len(globalStack.Top().nu))
		copy(agent.auth[h], globalStack.Top().nu)
	}

	globalStack.Update(1, agent.keyItr)
	agent.root = make([]byte, len(globalStack.Top().nu))
	copy(agent.root, globalStack.Top().nu)

	// don't forget to reset the key iterator
	agent.keyItr.Init(integatedSeed)

	return agent, nil
}

// refreshAuth updates auth path for next use
func (agent *MerkleAgent) refreshAuth() {
	//nextLeaf := agent.NumLeafUsed + 1
	nextLeaf := agent.keyItr.Offset()
	for h := uint32(0); h < agent.H; h++ {
		pow2Toh := uint32(1 << h)
		// nextLeaf % 2^h == 0
		if 0 == nextLeaf&(pow2Toh-1) {
			copy(agent.auth[h], agent.treeHashStacks[h].Top().nu)
			startingLeaf := (nextLeaf + pow2Toh) ^ pow2Toh
			agent.treeHashStacks[h].Init(startingLeaf, h)
			agent.stackIterators[h] = wots.NewKeyIterator(agent.GetSeedN(startingLeaf))
		}
	}
}

// refreshTreeHashStacks updates stack for next use
func (agent *MerkleAgent) refreshTreeHashStacks() {
	numOp := 2*agent.H - 1
	for i := uint32(0); i < numOp; i++ {
		globalLowest := uint32(math.MaxUint32)
		var focus uint32
		for h := uint32(0); h < agent.H; h++ {
			localLowest := agent.treeHashStacks[h].LowestTailHeight()
			if localLowest < globalLowest {
				globalLowest = localLowest
				focus = h
			}
		}
		agent.treeHashStacks[focus].Update(1, agent.stackIterators[focus])
	}
}

// Traverse updates both auth path and retained stack for next use
func (agent *MerkleAgent) Traverse() {
	agent.refreshAuth()
	agent.refreshTreeHashStacks()
}

// MerkleSig is the container for the signature generated
//	according to MSS
type MerkleSig struct {
	Leaf   uint32
	LeafPk *wots.PublicKey
	WtnSig *wots.WinternitzSig
	Auth   [][]byte
}

// Sign produces a Merkle signature
func Sign(agent *MerkleAgent, hash []byte) (*wots.PrivateKey, *MerkleSig, error) {
	merkleSig := new(MerkleSig)
	merkleSig.Leaf = agent.keyItr.Offset()

	sk, err := agent.keyItr.Next()
	merkleSig.WtnSig, err = wots.Sign(sk, hash)
	if nil != err {
		return nil, nil, errors.New("unexpected error occurs during signing")
	}

	// fill in the public key deriving leaf
	merkleSig.LeafPk = &sk.PublicKey
	// copy the auth path
	merkleSig.Auth = make([][]byte, len(agent.auth))
	for i := range agent.auth {
		merkleSig.Auth[i] = make([]byte, len(agent.auth[i]))
		copy(merkleSig.Auth[i], agent.auth[i])
	}

	// update auth path
	agent.Traverse()

	return sk, merkleSig, nil
}

// Verify verifies a Merkle signature
func Verify(root []byte, hash []byte, merkleSig *MerkleSig) bool {
	if (nil == merkleSig) || (!wots.Verify(merkleSig.LeafPk, hash, merkleSig.WtnSig)) {
		return false
	}

	H := len(merkleSig.Auth)
	// index of node in current height h
	idx := merkleSig.Leaf
	hashFunc := config.HashFunc()

	//parentHash := wots.HashPk(merkleSig.LeafPk)
	parentHash := HashPk(merkleSig.LeafPk)
	for h := 0; h < H; h++ {
		hashFunc.Reset()
		if 1 == idx%2 { // idx is odd, i.e., a right node
			hashFunc.Write(merkleSig.Auth[h])
			hashFunc.Write(parentHash)
		} else {
			hashFunc.Write(parentHash)
			hashFunc.Write(merkleSig.Auth[h])
		}
		// level up
		parentHash = hashFunc.Sum(nil)
		idx = idx >> 1
	}

	return bytes.Equal(parentHash, root)
}

// return the verification root
func (agent *MerkleAgent) Root() []byte {
	return agent.root
}

// get the nth Seed
func (agent *MerkleAgent) GetSeedN(n uint32) []byte {
	rand := rand.New(agent.seed0)
	for i := uint32(0); i < n; i++ {
		rand.NextState()
	}
	return rand.ExportSeed()
}
