package main

import (
	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/chaindb_config"
	"github.com/shynur/chaindb/discovery"
	"github.com/shynur/chaindb/monitor"
	"github.com/shynur/chaindb/validator"
)

func main() {
	discovery.Start(chaindb_config.DiscoveryInterval)

	validator.StartBlockDeliveryServer()
	validator.StartTryPickBlocks(LocalChain)

	blockchain.StartMining(LocalChain, &LocalTxPool, validator.Propose)

	monitor.Start(LocalChain, &LocalTxPool)
	StartUserService()

	select {}
}
