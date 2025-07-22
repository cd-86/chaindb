package blockchain

import (
	"math/rand/v2"
	"slices"
	"time"

	chaindb_config "github.com/shynur/chaindb/config"
)

func StartMining(chain *Chain, tx_pool *TxPool) {
	const blk_time = chaindb_config.BlockTime
}

func BuildBlockThenAppend_sync(chain *Chain, tx_pool *TxPool) {
	candidates := func() map[uint32][]Transaction {
		tx_bat := tx_pool.ExchangeNil_sync()

		candidates := make(map[uint32][]Transaction)
		chain.lock.RLock()
		defer chain.lock.RUnlock()
		for _, tx := range tx_bat {
			if chain.maybeValidNewTx(tx) {
				candidates[tx.OwnerID] = append(candidates[tx.OwnerID], tx)
			}
		}

		return candidates
	}()

	for _, txs := range candidates {
		slices.SortFunc(
			txs,
			func(tx1, tx2 Transaction) int {
				return int(tx1.Nonce) - int(tx2.Nonce)
			},
		)
	}

	blk := Block{
		Timestamp: time.Duration(time.Now().UnixNano()).Seconds(),
		UUID:      rand.Uint32(),
		ParentUUID: func() Block {
			chain.lock.RLock()
			defer chain.lock.RUnlock()
			return chain.blocks[len(chain.blocks)-1]
		}().UUID,
	}
}
