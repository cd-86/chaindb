package blockchain

type Block struct {
	Timestamp         float64       `json:"Timestamp"`
	receivedTimestamp float64       // 为 0 表示是本机自己挖出来的区块.
	UUID              uint32        `json:"UUID"`
	ParentUUID        uint32        `json:"ParentUUID"`
	Transactions      []Transaction `json:"Transactions"`
}

func (blk Block) getSeenTimestamp() float64 {
	if blk.receivedTimestamp != 0 {
		return blk.receivedTimestamp
	} else {
		return blk.Timestamp
	}
}
