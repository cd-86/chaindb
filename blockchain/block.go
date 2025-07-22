package blockchain

import (
	"errors"
	"slices"
	"sync"
	"time"
)

type Block struct {
	Timestamp         float64 `json:"Timestamp"`
	receivedTimestamp float64 // 为 0 表示是本机自己挖出来的区块.

	UUID       uint32 `json:"UUID"`
	Height     uint32 `json:"Height"`
	ParentUUID uint32 `json:"ParentUUID,omitempty"`

	Transactions []Transaction `json:"Transactions,omitempty"`
}

var BlockCache = func() *sync.Map {
	var cache sync.Map

	the_genesis_block := Block{
		Timestamp: time.Duration(time.Now().UnixNano()).Seconds(),
		UUID:      0,
	}

	cache.Store(the_genesis_block.UUID, the_genesis_block)

	return &cache
}() // UUID -> Block

func (blk Block) getSeenTimestamp() float64 {
	if blk.receivedTimestamp != 0 {
		return blk.receivedTimestamp
	} else {
		return blk.Timestamp
	}
}

func (tail Block) nextNonceOf(owner uint32) uint32 {
	for blk := tail; blk.Height != 0; func() Block {
		blk, _ := BlockCache.Load(blk.ParentUUID)
		return blk.(Block)
	}() {
		for _, tx_in_fork := range slices.Backward(blk.Transactions) {
			if tx_in_fork.OwnerID == owner {
				return tx_in_fork.Nonce + 1
			}
		}
	}
	return 0
}

func (blk Block) Previous() (previous Block, err error) {
	if blk.UUID == 0 {
		err = errors.New("创世区块没有父区块")
		return
	}
	previous_, _ := BlockCache.Load(blk.ParentUUID)
	previous = previous_.(Block)
	return
}
