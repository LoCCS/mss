package mss

import (
	"errors"
	"math"

	wots "github.com/sammy00/mss/ots/winternitz"
)

// MerkleSS implements the Merkle signature scheme
type MerkleSS struct {
	H              uint32
	NumLeafUsed    uint32
	auth           [][]byte
	root           []byte
	keyItr         *wots.SkPkIterator
	treeHashStacks []*TreeHashStack
}

// NewMerkleSS makes a fresh Merkle signing routine
//	by running the generate key and setup procedure
func NewMerkleSS(H uint32, seed []byte) (*MerkleSS, error) {
	if H < 2 {
		return nil, errors.New("H should be larger than 1")
	}

	merkleSS := new(MerkleSS)
	merkleSS.H = H
	merkleSS.auth = make([][]byte, H)
	merkleSS.keyItr = wots.NewSkPkIterator(seed)
	merkleSS.treeHashStacks = make([]*TreeHashStack, H)

	globalStack := NewTreeHashStack(0, H+1)
	//numLeaf := ((1 << H) - 1)

	for h := uint32(0); h < H; h++ {
		globalStack.Update(1, merkleSS.keyItr)

		merkleSS.treeHashStacks[h] = NewTreeHashStack(0, h)
		merkleSS.treeHashStacks[h].nodeStack.Push(globalStack.Top())

		globalStack.Update((1<<(h+1))-1, merkleSS.keyItr)
		merkleSS.auth[h] = make([]byte, len(globalStack.Top().nu))
		copy(merkleSS.auth[h], globalStack.Top().nu)
	}

	globalStack.Update(1, merkleSS.keyItr)
	merkleSS.root = make([]byte, len(globalStack.Top().nu))
	copy(merkleSS.root, globalStack.Top().nu)

	return merkleSS, nil
}

func (merkleSS *MerkleSS) refreshAuth() {
	nextLeaf := merkleSS.NumLeafUsed + 1
	for h := uint32(0); h < merkleSS.H; h++ {
		pow2Toh := uint32(1 << h)
		// nextLeaf % 2^h == 0
		if 0 == nextLeaf&pow2Toh {
			copy(merkleSS.auth[h], merkleSS.treeHashStacks[h].Top().nu)
			startingLeaf := (nextLeaf + pow2Toh) ^ pow2Toh
			merkleSS.treeHashStacks[h].Init(startingLeaf, h)
		}
	}
}
func (merkleSS *MerkleSS) refreshTreeHashStacks() {
	numOp := 2*merkleSS.H - 1
	for i := uint32(0); i < numOp; i++ {
		globalLowest := uint32(math.MaxUint32)
		var focus uint32
		for h := uint32(0); h < merkleSS.H; h++ {
			localLowest := merkleSS.treeHashStacks[h].LowestTailHeight()
			if localLowest < globalLowest {
				globalLowest = localLowest
				focus = h
			}
		}
		merkleSS.treeHashStacks[focus].Update(1, merkleSS.keyItr)
	}
}
func (merkleSS *MerkleSS) Traverse() {
	merkleSS.refreshAuth()
	merkleSS.refreshTreeHashStacks()
	merkleSS.NumLeafUsed++
}
