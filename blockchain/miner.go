package blockchain

import (
	"log"
	"math/rand/v2"
	"sync/atomic"
	"time"

	chaindb_config "github.com/shynur/chaindb/config"
	"github.com/shynur/chaindb/discovery"
)

func StartMining(chain *Chain, tx_pool *TxPool) {
	go func() {
		const check_interval = chaindb_config.BlockTime / 2
		for time.Sleep(check_interval); ; time.Sleep(check_interval) {
			// 整个区块链网络在某个时间节点生成新区块的概率:
			// p = 1 - 1/(1 + 当前block_time/BlockTime), p∈[0, 1)
			// 当前block_time 增加, p 增加, 使 当前block_time 减小.
			// 当前block_time = 预期时, p = 50%.
			// 因此每半个 BlockTime 检查一次, 也就是 check_interval.
			// 这样两个 block 之间的时间间隔的数学期望就是一个 BlockTime.
			p := 1 - 1/
				(1+
					(float64(chain.AvgBlockTime(chaindb_config.DifficultyAdjustmentWindow))/
						float64(chaindb_config.BlockTime)))

			// 当前节点生成新区块的概率:
			p /= float64(len(*discovery.ActiveNodes.Load())) + 1

			log.Printf("本机在当前检查点产生区块的概率为: %.2f%%\n", p*100)
			if rand.Float64() < p {
				log.Println("开采新区块...")
				BuildBlockTryAppend(chain, tx_pool)
			}
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
	log.Printf("已开采新区块, UUID=%d, Height=%d\n", blk.UUID, blk.Height)

	chain.TrySwitchHead(blk.UUID)
}
