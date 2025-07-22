package blockchain

import (
	"math/rand/v2"
	"sync/atomic"
	"time"

	chaindb_config "github.com/shynur/chaindb/config"
	"github.com/shynur/chaindb/discovery"
)

func StartMining(chain *Chain, tx_pool *TxPool) {
	go func() {
		for {
			// 整个网络下一秒生成新区块的概率:
			p := float64(
				chain.AvgBlockTime(chaindb_config.DifficultyAdjustmentWindow),
			) / float64(chaindb_config.BlockTime)

			// 当前节点生成新区块的概率:
			p /= float64(len(*discovery.ActiveNodes.Load())) + 1

			if rand.Float64() < p {
				BuildBlockTryAppend(chain, tx_pool)
			}

			time.Sleep(chaindb_config.BlockTime)
		}
	}()
}

func BuildBlockTryAppend(chain *Chain, tx_pool *TxPool) {
	blk := ForkFrom((*atomic.Uint32)(chain).Load())

	txin := make(chan Transaction)
	go tx_pool.PopFunc(
		txin,
		func(tx Transaction) string {
			switch {
			// Tx 太新了.
			case tx.Nonce > blk.nextNonceOf(tx.OwnerID):
				return "skip"
			// Tx 刚好是下一个.
			case tx.Nonce == blk.nextNonceOf(tx.OwnerID):
				return "pop"
			// Tx 太旧了.
			default:
				// 太久了也是要弹出的, 只是我们会直接把它扔掉罢了.
				return "discard"
			}
		},
	)
	for tx := range txin {
		blk.Transactions = append(blk.Transactions, tx)
	}

	BlockCache.Store(blk.UUID, blk)
	chain.TrySwitchHead(blk.UUID)
}
