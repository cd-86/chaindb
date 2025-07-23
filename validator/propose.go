package validator

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/chaindb_config"
	"github.com/shynur/chaindb/discovery"
)

func StartProposePeriodically(chain *blockchain.Chain, interval time.Duration) {
	last_proposed := func() blockchain.Block {
		the_genesis_blk, _ := blockchain.BlockCache.Load(uint32(0))
		return the_genesis_blk.(blockchain.Block)
	}()
	go func() {
		for ; ; time.Sleep(interval) {
			head := chain.Head()
			if head.MinerAddress != chaindb_config.MinerAddress {
				continue
			}
			if head.UUID == last_proposed.UUID {
				continue
			}

			proposeBlock(head)
			last_proposed = head
		}
	}()
}

func proposeBlock(blk blockchain.Block) {
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
