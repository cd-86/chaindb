package validator

import (
	"log"
	"net"
	"os"
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
		max_blk_bufsize := os.Getpagesize()
		for {
			blk_obj := make([]byte, max_blk_bufsize)
			n, _, _ := conn.ReadFrom(blk_obj)

			if n >= len(blk_obj) {
				log.Printf("[  Pull ] 接收到的区块数据过大, 被丢弃!\n")
				max_blk_bufsize *= 2
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
