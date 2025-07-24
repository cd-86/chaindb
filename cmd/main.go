package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/chaindb_config"
	"github.com/shynur/chaindb/discovery"
	"github.com/shynur/chaindb/validator"
)

func main() {
	discovery.Start(chaindb_config.DiscoveryInterval)

	validator.StartBlockDeliveryServer()
	validator.StartTryPickBlocks(LocalChain)
	blockchain.StartMining(LocalChain, &LocalTxPool, validator.Propose)

	for i := 0; i != 10_0000; i++ {
		tx := blockchain.Transaction{
			OwnerID: rand.Uint32N(20000),
			Nonce:   0,
			Data:    "000000000000000000000000000",
		}
		LocalTxPool.Add(tx)

		time.Sleep(time.Millisecond)
	}

	fmt.Println(LocalChain)
	select {}
}
