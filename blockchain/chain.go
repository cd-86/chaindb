package blockchain

import (
	"log"
	"math/rand/v2"
	"sync/atomic"
	"time"

	chaindb_config "github.com/shynur/chaindb/config"
)

type Chain atomic.Pointer[Block]

func New() *Chain {
	the_genesis_block := Block{
		Timestamp: time.Duration(time.Now().UnixNano()).Seconds(),
	}
	var chain atomic.Pointer[Block]
	chain.Store(&the_genesis_block)
	return (*Chain)(&chain)
}

func ForkFrom(parent *Block) Block {
	if parent == nil {
		log.Fatalln("不允许从不存在的区块分叉")
	}
	return Block{
		UUID:       rand.Uint32(),
		Height:     parent.Height + 1,
		ParentUUID: parent.ParentUUID,
		parent:     parent,

		Timestamp: time.Duration(time.Now().UnixNano()).Seconds(),
	}
}

func (chain *Chain) Fork() Block {
	return ForkFrom((*atomic.Pointer[Block])(chain).Load())
}

func (chain *Chain) AvgBlockTime(samples uint) time.Duration {
	if samples <= 1 {
		return chaindb_config.BlockTime
	}

	tail := (*atomic.Pointer[Block])(chain).Load()
	num_blocks := tail.Height + 1

	if num_blocks == 1 {
		return chaindb_config.BlockTime
	}

	if num_blocks <= uint32(samples) {
		the_genesis_block := tail
		for the_genesis_block.UUID != 0 {
			the_genesis_block = the_genesis_block.parent
		}
		return time.Duration(
			(tail.getSeenTimestamp() - the_genesis_block.getSeenTimestamp()) /
				(float64(samples - 1)),
		)
	}

	begin := tail
	for range samples - 1 {
		begin = begin.parent
	}
	return time.Duration(
		(tail.getSeenTimestamp() - begin.getSeenTimestamp()) /
			(float64(samples - 1)),
	)
}
