package validator

import "github.com/shynur/chaindb/blockchain"

func StartBlockDeliveryServer() {
	_ = blockchain.BlockCache
}
