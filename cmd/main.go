package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/shynur/chaindb/blockchain"
	chaindb_config "github.com/shynur/chaindb/config"
	"github.com/shynur/chaindb/discovery"
)

func main() {
	discovery.Start(chaindb_config.DiscoveryInterval)

	blockchain.StartMining(LocalChain, &LocalTxPool)

	for i := 0; i != 1000; i++ {
		tx := blockchain.Transaction{
			OwnerID: rand.Uint32N(5),
			Nonce:   rand.Uint32N(5),
		}
		LocalTxPool.Add(tx)

		time.Sleep(30 * time.Millisecond)
	}

	fmt.Printf("%+v", LocalChain)
}
