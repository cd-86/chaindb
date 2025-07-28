package validator

import (
	"log"
	"net"
	"strconv"

	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/chaindb_config"
	"github.com/shynur/chaindb/discovery"
)

func Propose(blk blockchain.Block) {
	blk_obj := blk.ToGob()
	log.Printf("[Propose] 即将被广播的区块: size=%dB\n", len(blk_obj))

	for _, host := range *discovery.ActiveNodes.Load() {
		go func() {
			conn, _ := net.Dial(
				"udp",
				net.JoinHostPort(
					host,
					strconv.Itoa(chaindb_config.UDPPortPickBlock),
				),
			)
			defer conn.Close()
			conn.Write(blk_obj)
		}()
	}
}
