package main

import (
	"math"
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

	go func() {
		for ; ; time.Sleep(15 * time.Second) {
			//log.Println("[ Chain ] 可见的全量最长区块链 ", LocalChain)
		}
	}()
	const loop_cnt = 100_0000
	for range loop_cnt {
		const num_owners = 20
		tx := blockchain.Transaction{
			OwnerID: rand.Uint32N(num_owners),
			Nonce: rand.Uint32N(
				1 + uint32(math.Sqrt(loop_cnt/num_owners)),
			),
		}
		LocalTxPool.Add(tx)
		time.Sleep(100 * time.Millisecond)
	}
}
