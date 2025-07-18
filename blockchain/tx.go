package blockchain

import "container/list"

type Transaction struct {
	OwnerID uint32 `json:"OwnerID"`
	Nonce   uint32 `json:"Nonce"`
	Data    string `json:"Data"`
}

var TxQueue list.List
