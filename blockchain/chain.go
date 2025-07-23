package blockchain

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"sync/atomic"
	"time"

	"github.com/shynur/chaindb/chaindb_config"
)

// 区块链最后一个区块的 UUID.
type Chain atomic.Uint32

func (chain *Chain) Head() Block {
	head_uuid := (*atomic.Uint32)(chain).Load()
	head, _ := BlockCache.Load(head_uuid)
	return head.(Block)
}

func (chain *Chain) String() string {
	reversed_chain := []Block{chain.Head()}
	for reversed_chain[len(reversed_chain)-1].UUID != 0 {
		previous, _ := reversed_chain[len(reversed_chain)-1].Previous()
		reversed_chain = append(
			reversed_chain,
			previous,
		)
	}
	slices.Reverse(reversed_chain)

	return fmt.Sprintf("%+v", reversed_chain)
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

// 查看最近两个 locally mined block 区间的平均区块时间.
// 只对 locally mined block 进行采样,
// 因为其它 block 的 Timestamp 并不是相对于本机的时间线.
func (chain *Chain) AvgBlockTime() time.Duration {
	right := chain.Head()

	for right.MinerAddress != chaindb_config.MinerAddress {
		right, _ = right.Previous()
	}

	if right.Height == 0 {
		return chaindb_config.BlockTime
	}

	left, _ := right.Previous()
	for left.MinerAddress != chaindb_config.MinerAddress {
		left, _ = left.Previous()
	}

	return time.Duration(
		(right.Timestamp - left.Timestamp) /
			float64((right.Height - left.Height)) *
			1e9,
	)
}

func (chain *Chain) TrySwitchHead(new_tail_uuid uint32) (switched bool) {
	new_tail, _ := BlockCache.Load(new_tail_uuid)

	if new_tail.(Block).Height > chain.Head().Height {
		(*atomic.Uint32)(chain).Store(new_tail_uuid)
		return true
	}
	return false
}
