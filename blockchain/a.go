package blockchain

import "sync"

type Transaction struct {
	OwnerID uint32 `json:"OwnerID"`
	Nonce   uint32 `json:"Nonce"`
	Data    string `json:"Data"`
}

type Block struct {
	Timestamp    float64       `json:"Timestamp"`
	UUID         uint32        `json:"UUID"`
	ParentUUID   uint32        `json:"ParentUUID"`
	Transactions []Transaction `json:"Transactions"`
}

type Chain struct {
	blocks []Block
	lock   sync.RWMutex
}

func NewChain() *Chain {
	var chain Chain
	chain.blocks = []Block{
		//TODO {Timestamp: }
	}
	return &chain
}
