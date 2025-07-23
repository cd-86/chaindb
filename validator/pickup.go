package validator

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"strconv"

	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/chaindb_config"
)

func StartTryPickBlocks() {
	conn, err := net.ListenPacket(
		"udp",
		net.JoinHostPort(
			"",
			strconv.Itoa(chaindb_config.UDPPortPickBlock),
		),
	)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			blk_obj := make([]byte, 65536)
			n, _, _ := conn.ReadFrom(blk_obj)

			if n >= len(blk_obj) {
				log.Printf("接收到的区块数据过大, 被丢弃!\n")
				continue
			} else {
				log.Printf("接收到其它矿工开采的区块: size=%dB\n", n)
			}

			blk_obj = blk_obj[:n]

			var blk blockchain.Block
			gob.NewDecoder(bytes.NewBuffer(blk_obj)).Decode(&blk)
		}
	}()
}
