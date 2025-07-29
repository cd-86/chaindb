package monitor

import (
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/chaindb_config"
)

var (
	theChain  *blockchain.Chain
	theTxPool *blockchain.TxPool
)

func Start(chain *blockchain.Chain, tx_pool *blockchain.TxPool) {
	theChain, theTxPool = chain, tx_pool

	server := http.NewServeMux()

	registerTxService(server)
	registerBlockService(server)
	registerChainService(server)
	registerNodeService(server)

	go func() {
		err := http.ListenAndServe(
			net.JoinHostPort(
				"",
				strconv.Itoa(chaindb_config.TCPPortMonitor),
			),
			server,
		)
		log.Printf("[Monitor] Error: %v\n", err)
	}()
}
