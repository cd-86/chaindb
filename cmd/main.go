package main

import (
	"log"
	"runtime"
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
	StartUserService()

	for ; ; time.Sleep(chaindb_config.BlockTime / 2) {
		log.Println("[-] 当前 goroutine 的数量", runtime.NumGoroutine())
	}
	// select {}
}
