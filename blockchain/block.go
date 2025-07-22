package blockchain

import "slices"

type Block struct {
	Timestamp         float64 `json:"Timestamp"`
	receivedTimestamp float64 // 为 0 表示是本机自己挖出来的区块.

	UUID       uint32 `json:"UUID"`
	Height     uint32 `json:"Height"`
	ParentUUID uint32 `json:"ParentUUID,omitempty"`
	parent     *Block

	Transactions []Transaction `json:"Transactions,omitempty"`
}

func (blk Block) getSeenTimestamp() float64 {
	if blk.receivedTimestamp != 0 {
		return blk.receivedTimestamp
	} else {
		return blk.Timestamp
	}
}

func (tail Block) nextNonceOf(owner uint32) uint32 {
	for blk := &tail; blk != nil; blk = tail.parent {
		for _, tx_in_fork := range slices.Backward(blk.Transactions) {
			if tx_in_fork.OwnerID == owner {
				return tx_in_fork.Nonce + 1
			}
		}
	}
	return 0
}
