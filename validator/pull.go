package validator

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/chaindb_config"
	"github.com/shynur/chaindb/discovery"
	"github.com/shynur/chaindb/validator/block_cdn"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func pullOneBlock(blk_uuid uint32) (blockchain.Block, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	block_request := make(chan []byte)

	var requests sync.WaitGroup
	for _, ip := range *discovery.ActiveNodes.Load() {
		requests.Add(1)
		go func() {
			defer requests.Done()

			conn, err := grpc.NewClient(
				net.JoinHostPort(
					ip.String(),
					strconv.Itoa(chaindb_config.TCPPortBlockPCDN),
				),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				return
			} else {
				defer conn.Close()
			}

			client := block_cdn.NewBlockDeliveryClient(conn)

			resp, err := client.GetBlock(
				ctx,
				&block_cdn.BlockUUID{UUID: blk_uuid},
			)
			if err != nil {
				return
			}

			select {
			case block_request <- resp.Buffer:
				cancel()
			default:
			}
		}()
	}

	select {
	case buf := <-block_request:
		fmt.Printf("Got result: %v\n", buf)
	case <-time.After(5 * time.Second):
		fmt.Fprintln(os.Stderr, "timeout")
	}

	requests.Wait()
}

// head_uuid 块已经存在, 然后从 head_uuid 开始拉取.
// 如果够长, 则 switch 到该 head_uuid.
func pull(chain *blockchain.Chain, head_uuid uint32) {
	head, _ := blockchain.BlockCache.Load(head_uuid)
	previous_uuid := head.(blockchain.Block).ParentUUID

	for {
		_, exist := blockchain.BlockCache.Load(previous_uuid)
		if exist {
			break
		}

		previous_blk, err := pullOneBlock(previous_uuid)
		if err != nil {
			return
		}

		blockchain.BlockCache.Store(previous_uuid, previous_blk)
		previous_uuid = previous_blk.ParentUUID
	}

	chain.TrySwitchHead(head_uuid)
}
