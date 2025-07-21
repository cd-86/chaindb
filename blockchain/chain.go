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

func (chain *Chain) GetOwners() (
	owners []struct {
		OwnerID           uint32 `json:"OwnerID"`
		ConfirmationScore uint32 `json:"ConfirmationScore"`
	},
) {
	chain.lock.RLock()
	defer chain.lock.RUnlock()

	for i, block := range chain.blocks {
		for _, tx := range block.Transactions {
			if tx.Nonce == 0 {
				owners = append(
					owners,
					struct {
						OwnerID           uint32 `json:"OwnerID"`
						ConfirmationScore uint32 `json:"ConfirmationScore"`
					}{
						OwnerID:           tx.OwnerID,
						ConfirmationScore: uint32(len(chain.blocks) - i - 1),
					},
				)
			}
		}
	}

	return owners
}

func (chain *Chain) GetTxOwnedBy(owner_id uint32) (
	transactions []struct {
		ConfirmationScore uint32 `json:"ConfirmationScore"`
		Data              string `json:"Data"`
	},
) {
	chain.lock.RLock()
	defer chain.lock.RUnlock()

	for i, block := range chain.blocks {
		for _, tx := range block.Transactions {
			if tx.OwnerID == owner_id {
				transactions = append(
					transactions,
					struct {
						ConfirmationScore uint32 `json:"ConfirmationScore"`
						Data              string `json:"Data"`
					}{
						ConfirmationScore: uint32(len(chain.blocks) - i - 1),
						Data:              tx.Data,
					},
				)
			}
		}
	}

	return transactions
}
