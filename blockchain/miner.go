package blockchain

import (
	"sync/atomic"

	chaindb_config "github.com/shynur/chaindb/config"
)

func StartMining(chain *Chain, tx_pool *TxPool) {
	const blk_time = chaindb_config.BlockTime
}

func BuildBlockTryAppend(chain *Chain, tx_pool *TxPool) {
	blk := ForkFrom((*atomic.Uint32)(chain).Load())

	txin := make(chan Transaction)
	go tx_pool.PopFunc(
		txin,
		func(tx Transaction) bool {
			blk.nextNonceOf(tx.OwnerID)
			return true
		},
	)
	for tx := range txin {
		blk.Transactions = append(blk.Transactions, tx)
	}

	BlockCache.Store(blk.UUID, blk)
	chain.TrySwitchHead(blk.UUID)
}
