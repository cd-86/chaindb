package blockchain

import (
	"slices"

	chaindb_config "github.com/shynur/chaindb/config"
)

func StartMining(chain *Chain, tx_pool *TxPool) {
	const blk_time = chaindb_config.BlockTime
}

func BuildBlockThenAppend_sync(chain *Chain, tx_pool *TxPool) {
	candidates := func() map[uint32][]Transaction {
		chain.lock.Lock()
		defer chain.lock.Unlock()
		return func() (candidates map[uint32][]Transaction) {
			tx_pool.lock.Lock()
			defer tx_pool.lock.Unlock()
			for tx := tx_pool.transactions.Front(); tx != nil; tx = tx.Next() {
				if chain.maybeValidNewTx(tx.Value.(Transaction)) {
					candidates[tx.Value.(Transaction).OwnerID] = append(
						candidates[tx.Value.(Transaction).OwnerID],
						tx.Value.(Transaction),
					)
					tx_pool.transactions.Remove(tx)
				}
			}
			return
		}()
	}()

	for _, txs := range candidates {
		slices.SortFunc(
			txs,
			func(tx1, tx2 Transaction) int {
				return int(tx1.Nonce) - int(tx2.Nonce)
			},
		)
	}
	/*
		blk := Block{
			Timestamp: time.Duration(time.Now().UnixNano()).Seconds(),
			UUID:      rand.Uint32(),
			ParentUUID: func() Block {
				chain.lock.RLock()
				defer chain.lock.RUnlock()
				return chain.blocks[len(chain.blocks)-1]
			}().UUID,
		}
	*/
}
