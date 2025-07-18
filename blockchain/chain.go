package blockchain

import "sync"

type Chain struct {
	blocks []Block
	lock   sync.RWMutex
}

func NewChain() *Chain {
	var chain Chain
	chain.blocks = []Block{
		//TODO {Timestamp: }
	}
	return &chain
}
