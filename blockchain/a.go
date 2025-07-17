package blockchain

type Transaction struct{}

type Block struct {
	Timestamp  float64 `json:"Timestamp"`
	UUID       uint32  `json:"UUID"`
	ParentUUID uint32  `json:"ParentUUID"`
}
