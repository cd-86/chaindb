package validator

import "github.com/shynur/chaindb/blockchain"

// blk_uuid 块已经存在, 然后从 blk_uuid 开始拉取.
// 如果够长, 则 switch 到该 blk_uuid.
func pull(chain *blockchain.Chain, blk_uuid uint32) {
	//
}
