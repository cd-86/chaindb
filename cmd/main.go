package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/chaindb_config"
	"github.com/shynur/chaindb/discovery"
)

func main() {
	discovery.Start(chaindb_config.DiscoveryInterval)

	blockchain.StartMining(LocalChain, &LocalTxPool)

	for i := 0; i != 10_0000; i++ {
		tx := blockchain.Transaction{
			OwnerID: rand.Uint32N(20),
			Nonce:   rand.Uint32N(10),
		}
		LocalTxPool.Add(tx)

		time.Sleep(time.Millisecond)
	}

	fmt.Println(LocalChain)
}
