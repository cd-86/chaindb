package main

import (
	"fmt"
	"log"
	"time"

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

	for ; ; time.Sleep(15 * time.Second) {
		var blocks string
		for _, blk := range LocalChain.Snapshot() {
			blocks += fmt.Sprintf("%+v\n", blk)
		}
		log.Printf("[ Chain ] 可见的全量最长区块链:\n%s\n", blocks)
	}
}
