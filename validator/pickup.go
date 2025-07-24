package validator

import (
	"log"
	"net"
	"strconv"

	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/chaindb_config"
)

func StartTryPickBlocks(chain *blockchain.Chain) {
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

			var blk blockchain.Block
			blk.FromGob(blk_obj[:n])
			blockchain.BlockCache.Store(blk.UUID, blk)

			if blk.Height > chain.Head().Height {
				go pull(chain, blk.UUID)
			}
		}
	}()
}
