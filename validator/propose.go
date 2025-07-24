package validator

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"strconv"

	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/chaindb_config"
	"github.com/shynur/chaindb/discovery"
)

func Propose(blk blockchain.Block) {
	var blk_obj bytes.Buffer
	gob.NewEncoder(&blk_obj).Encode(blk)
	log.Printf("即将被广播的区块: size=%dB\n", blk_obj.Len())

	for _, ip_addr := range *discovery.ActiveNodes.Load() {
		go func() {
			conn, _ := net.Dial(
				"udp",
				net.JoinHostPort(
					ip_addr.String(),
					strconv.Itoa(chaindb_config.UDPPortPickBlock),
				),
			)
			defer conn.Close()
			conn.Write(blk_obj.Bytes())
		}()
	}
}
