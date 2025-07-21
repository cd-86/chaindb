package blockchain

import (
	"sync"
)

type Chain struct {
	blocks []Block
	lock   sync.RWMutex
}

func New() *Chain {
	var the_genesis_block Block
	return &Chain{
		blocks: []Block{the_genesis_block},
	}
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

func (chain *Chain) GetTxOwnedBy(owner_id uint32) {
	type tx_t = struct {
		ConfirmationScore uint32 `json:"ConfirmationScore"`
		Data              string `json:"Data"`
	}
	transactions := []tx_t{}

	chain.lock.RLock()
	defer chain.lock.RUnlock()
}
