package blockchain

import (
	chaindb_config "github.com/shynur/chaindb/config"
)

func StartMining(chain *Chain, tx_pool *TxPool) {
	const blk_time = chaindb_config.BlockTime
}

func BuildBlockThenAppend_sync(chain *Chain, tx_pool *TxPool) {
	chain.lock.Lock()
	defer chain.lock.Unlock()

	tx_candidates := []Transaction{}
	func() {
		tx_pool.lock.Lock()
		defer tx_pool.lock.Unlock()
		for tx := tx_pool.transactions.Front(); tx != nil; tx = tx.Next() {
			if chain.maybeValidNewTx(tx.Value.(Transaction)) {
				tx_candidates = append(tx_candidates, tx.Value.(Transaction))
				tx_pool.transactions.Remove(tx)
			}
		}
	}()
}
