package blockchain

import (
	"sync"
)

type Chain struct {
	blocks []Block
	lock   sync.RWMutex
}

var LocalChain = Chain{
	blocks: []Block{
		{ /* 创世区块 */ },
	},
}

func (chain *Chain) GetOwners() []struct {
	OwnerID           uint32 `json:"OwnerID"`
	ConfirmationScore uint32 `json:"ConfirmationScore"`
} {
	type owner_t = struct {
		OwnerID           uint32 `json:"OwnerID"`
		ConfirmationScore uint32 `json:"ConfirmationScore"`
	}
	owners := []owner_t{}

	chain.lock.RLock()
	defer chain.lock.RUnlock()

	for i, block := range chain.blocks {
		for _, tx := range block.Transactions {
			if tx.Nonce == 0 {
				owners = append(
					owners,
					owner_t{
						OwnerID:           tx.OwnerID,
						ConfirmationScore: uint32(len(chain.blocks) - i - 1),
					},
				)
			}
		}
	}

	return owners
}
