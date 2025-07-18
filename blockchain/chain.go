package blockchain

import "sync"

type Chain struct {
	blocks []Block
	lock   sync.RWMutex
}

var LocalChain = Chain{
	blocks: []Block{
		{ /* 创世区块 */ },
	},
}
