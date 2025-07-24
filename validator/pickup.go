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
				log.Printf("[  Pull ] 接收到的区块数据过大, 被丢弃!\n")
				continue
			} else {
				log.Printf("[  Pull ] 接收到其它矿工开采的区块: size=%dB\n", n)
			}

			var head blockchain.Block
			head.FromGob(blk_obj[:n])
			blockCacheDetached.Store(head.UUID, head)

			if head.Height > chain.Head().Height {
				go pull(chain, head)
			}
		}
	}()
}
