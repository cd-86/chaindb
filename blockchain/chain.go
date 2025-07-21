package blockchain

import (
	"log"
	"sync"
	"time"

	chaindb_config "github.com/shynur/chaindb/config"
)

type Chain struct {
	blocks []Block
	lock   sync.RWMutex
}

func New() *Chain {
	the_genesis_block := Block{
		Timestamp: time.Duration(time.Now().UnixNano()).Seconds(),
	}
	return &Chain{
		blocks: []Block{
			the_genesis_block,
		},
	}
}

// 有可能被加入到新区块所记录的交易列表中.
func (chain *Chain) validNewTx(tx Transaction) bool {
	return tx.Nonce >= uint32(len(chain.getTxOwnedBy(tx.OwnerID)))
}

func (chain *Chain) avgBlkTime() time.Duration {
	if num_blocks := len(chain.blocks); num_blocks == 0 {
		log.Fatalln("区块链的长度应当永远是正数才对, 默认有创世区块")
		return 0 // stupid gc
	} else if num_blocks == 1 {
		return chaindb_config.BlockTime
	} else if num_blocks <= 6 {
		return time.Duration(
			(chain.blocks[num_blocks-1].getSeenTimestamp() -
				chain.blocks[0].getSeenTimestamp()) /
				float64(num_blocks-1),
		)
	} else {
		return time.Duration(
			(chain.blocks[num_blocks-1].getSeenTimestamp() -
				chain.blocks[num_blocks-6].getSeenTimestamp()) /
				5.0,
		)
	}
}

func (chain *Chain) GetOwners_sync() (
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

func (chain *Chain) GetTxOwnedBy_sync(owner_id uint32) (
	transactions []struct {
		ConfirmationScore uint32 `json:"ConfirmationScore"`
		Data              string `json:"Data"`
	},
) {
	chain.lock.RLock()
	defer chain.lock.RUnlock()
	return chain.getTxOwnedBy(owner_id)
}

func (chain *Chain) getTxOwnedBy(owner_id uint32) (
	transactions []struct {
		ConfirmationScore uint32 `json:"ConfirmationScore"`
		Data              string `json:"Data"`
	},
) {
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
