package main

import (
	"math/rand/v2"
	"time"

	"github.com/shynur/chaindb/blockchain"
	chaindb_config "github.com/shynur/chaindb/config"
	"github.com/shynur/chaindb/discovery"
)

func main() {
	discovery.Start(chaindb_config.DiscoveryInterval)
	blockchain.StartMining(LocalChain, &LocalTxPool)

	for ; ; time.Sleep(100 * time.Millisecond) {
		tx := blockchain.Transaction{
			OwnerID: rand.Uint32N(5),
			Nonce:   rand.Uint32N(10),
		}
		LocalTxPool.Add(tx)
	}
}
