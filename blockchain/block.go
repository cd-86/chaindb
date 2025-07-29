package blockchain

import (
	"bytes"
	"encoding/gob"
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/shynur/chaindb/chaindb_config"
)

type Block struct {
	Timestamp    float64 `json:"Timestamp"`
	MinerAddress string  `json:"MinerAddress"`

	UUID       uint32 `json:"UUID"`
	Height     uint32 `json:"Height"`
	ParentUUID uint32 `json:"ParentUUID"`

	Transactions []Transaction `json:"Transactions,omitempty"`
}

var BlockCache = func() *sync.Map {
	var cache sync.Map

	the_genesis_block := Block{
		Timestamp:    time.Duration(time.Now().UnixNano()).Seconds(),
		MinerAddress: chaindb_config.MinerAddress,
		UUID:         0,
		Height:       0,
	}

	cache.Store(the_genesis_block.UUID, the_genesis_block)

	return &cache
}() // UUID:uint32 -> Block

func (tail Block) nextNonceOf(owner uint32) uint32 {
	for blk := tail; blk.Height != 0; blk, _ = blk.Previous() {
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

func (blk Block) ToGob() []byte {
	var obj bytes.Buffer
	gob.NewEncoder(&obj).Encode(blk)
	return obj.Bytes()
}
func (blk *Block) FromGob(obj []byte) error {
	err := gob.NewDecoder(bytes.NewBuffer(obj)).Decode(blk)
	return err
}
