package main

import (
	"github.com/shynur/chaindb/blockchain"
	chaindb_config "github.com/shynur/chaindb/config"
	"github.com/shynur/chaindb/discovery"
)

func main() {
	discovery.Start(chaindb_config.DiscoveryInterval)

	blockchain.StartMining(LocalChain, &LocalTxPool)
}
