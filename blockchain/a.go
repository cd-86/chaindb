package blockchain

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
