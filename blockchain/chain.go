package blockchain

import (
	"math/rand/v2"
	"slices"
	"sync/atomic"
	"time"

	"github.com/shynur/chaindb/chaindb_config"
)

// 区块链最后一个区块的 UUID.
type Chain atomic.Uint32

// 按插入顺序返回.
func (chain *Chain) ListOwners() (
	owners []struct {
		OwnerID           uint32
		ConfirmationScore uint32
	},
) {
	owners_earliest_tx := make(map[uint32]struct {
		ConfirmationScore uint32
		Age               uint32
	})
	head := chain.Head()
	num_txs := uint32(0)
	for blk := head; blk.UUID != 0; blk, _ = blk.Previous() {
		for _, tx := range func() []Transaction {
			transactions := slices.Clone(blk.Transactions)
			slices.Reverse(blk.Transactions)
			return transactions
		}() {
			owners_earliest_tx[tx.OwnerID] = struct {
				ConfirmationScore uint32
				Age               uint32
			}{
				ConfirmationScore: head.Height - blk.Height,
				Age:               num_txs,
			}

			num_txs++
		}
	}

	owners = []struct {
		OwnerID           uint32
		ConfirmationScore uint32
	}{}
	for owner, tx := range owners_earliest_tx {
		owners = append(
			owners,
			struct {
				OwnerID           uint32
				ConfirmationScore uint32
			}{
				OwnerID:           owner,
				ConfirmationScore: tx.ConfirmationScore,
			},
		)
	}
	slices.SortFunc(
		owners,
		func(owner1, owner2 struct {
			OwnerID           uint32
			ConfirmationScore uint32
		}) int {
			return -(int(owners_earliest_tx[owner1.OwnerID].Age) -
				int(owners_earliest_tx[owner2.OwnerID].Age))
		},
	)

	return
}

func (chain *Chain) Head() Block {
	head_uuid := (*atomic.Uint32)(chain).Load()
	head, _ := BlockCache.Load(head_uuid)
	return head.(Block)
}

func (chain *Chain) Snapshot() []Block {
	reversed_chain := []Block{chain.Head()}
	for reversed_chain[len(reversed_chain)-1].UUID != 0 {
		previous, _ := reversed_chain[len(reversed_chain)-1].Previous()
		reversed_chain = append(
			reversed_chain,
			previous,
		)
	}
	slices.Reverse(reversed_chain)

	return reversed_chain
}

func New() *Chain {
	var blk_uuid atomic.Uint32
	return (*Chain)(&blk_uuid)
}

func ForkFrom(parent_uuid uint32) Block {
	parent, _ := BlockCache.Load(parent_uuid)

	return Block{
		UUID:       rand.Uint32(),
		Height:     parent.(Block).Height + 1,
		ParentUUID: parent_uuid,

		MinerAddress: chaindb_config.MinerAddress,

		Timestamp: time.Duration(time.Now().UnixNano()).Seconds(),
	}
}

// 对最近两个 由同一个矿工挖出的 block 进行采样, 算出平均区块时间.
func (chain *Chain) AvgBlockTime() time.Duration {
	recently_mined := map[string]Block{}

	for blk := chain.Head(); blk.UUID != 0; blk, _ = blk.Previous() {
		if descendant_blk, exists := recently_mined[blk.MinerAddress]; exists {
			return time.Duration(
				(descendant_blk.Timestamp - blk.Timestamp) /
					float64((descendant_blk.Height - blk.Height)) *
					1e9,
			)
		}
		recently_mined[blk.MinerAddress] = blk
	}

	return chaindb_config.BlockTime
}

func (chain *Chain) TrySwitchHead(new_tail_uuid uint32) (switched bool) {
	new_tail, _ := BlockCache.Load(new_tail_uuid)

	if new_tail.(Block).Height > chain.Head().Height {
		(*atomic.Uint32)(chain).Store(new_tail_uuid)
		return true
	}
	return false
}
