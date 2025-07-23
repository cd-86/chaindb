package validator

import (
	"time"

	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/chaindb_config"
)

func StartProposePeriodically(chain *blockchain.Block, interval time.Duration) {
	go func() {
		for ; ; time.Sleep(chaindb_config.BlockTime / 2) {
		}
	}()
}
