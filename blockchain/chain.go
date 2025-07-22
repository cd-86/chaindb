package blockchain

import (
	"math/rand/v2"
	"sync/atomic"
	"time"

	chaindb_config "github.com/shynur/chaindb/config"
)

// 区块链最后一个区块的 UUID.
type Chain atomic.Uint32

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

		Timestamp: time.Duration(time.Now().UnixNano()).Seconds(),
	}
}

func (chain *Chain) AvgBlockTime(samples uint) time.Duration {
	if samples <= 1 {
		return chaindb_config.BlockTime
	}

	tail_, _ := BlockCache.Load((*atomic.Uint32)(chain).Load())
	tail := tail_.(Block)
	num_blocks := tail.Height + 1

	if num_blocks == 1 {
		return chaindb_config.BlockTime
	}

	if num_blocks <= uint32(samples) {
		the_genesis_block := tail
		for the_genesis_block.UUID != 0 {
			the_genesis_block, _ = the_genesis_block.Previous()
		}
		return time.Duration(
			(tail.getSeenTimestamp() - the_genesis_block.getSeenTimestamp()) /
				(float64(samples - 1)),
		)
	}

	begin := tail
	for range samples - 1 {
		begin, _ = begin.Previous()
	}
	return time.Duration(
		(tail.getSeenTimestamp() - begin.getSeenTimestamp()) /
			(float64(samples - 1)),
	)
}

func (chain *Chain) TrySwitchHead(new_tail_uuid uint32) (switched bool) {
	new_tail, _ := BlockCache.Load(new_tail_uuid)

	if new_tail.(Block).Height >= chain.Length() {
		(*atomic.Uint32)(chain).Store(new_tail_uuid)
		return true
	}
	return false
}

func (chain *Chain) Length() uint32 {
	tail, _ := BlockCache.Load((*atomic.Uint32)(chain).Load())
	return tail.(Block).Height + 1
}
