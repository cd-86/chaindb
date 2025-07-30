package validator

import (
	"context"
	"log"
	"net"
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

var blockCacheDetached sync.Map

func pullOneBlock(blk_uuid uint64) (blk blockchain.Block, err error) {
	cached_detached_block, exist := blockCacheDetached.Load(blk_uuid)
	if exist {
		return cached_detached_block.(blockchain.Block), nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	block_request := make(chan []byte, 1)

	for _, host := range *discovery.ActiveNodes.Load() {
		go func() {
			conn, err := grpc.NewClient(
				net.JoinHostPort(
					host,
					strconv.Itoa(chaindb_config.TCPPortBlockPCDN),
				),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				// 放开 gRPC 的所有限制:
				// rpc.WithInitialConnWindowSize(math.MaxInt32),
				// grpc.WithInitialWindowSize(math.MaxInt32),
				// grpc.WithMaxHeaderListSize(math.MaxUint32),
				// grpc.WithStaticConnWindowSize(math.MaxInt32),
				// grpc.WithStaticStreamWindowSize(math.MaxInt32),
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

	err = blk.FromGob(<-block_request)
	if err == nil {
		log.Printf(
			"[  Pull ] 拉取到区块 Block.UUID=%d, Block.Height=%d\n",
			blk.UUID,
			blk.Height,
		)
		blockCacheDetached.Store(blk_uuid, blk)
	}
	return
}

// head_uuid 块已经存在, 然后从 head_uuid 开始拉取.
// 如果够长, 则 switch 到该 head_uuid.
func pull(chain *blockchain.Chain, head blockchain.Block) {
	candidates := []blockchain.Block{head}

	previous_uuid := head.ParentUUID

	for {
		_, exist := blockchain.BlockCache.Load(previous_uuid)
		if exist {
			break
		}

		var previous_blk blockchain.Block
		select {
		case resp := <-func() <-chan struct {
			blockchain.Block
			error
		} {
			resp := make(chan struct {
				blockchain.Block
				error
			})
			go func() {
				previous_blk, err := pullOneBlock(previous_uuid)
				resp <- struct {
					blockchain.Block
					error
				}{previous_blk, err}
			}()
			return resp
		}():
			if resp.error != nil {
				return
			}
			previous_blk = resp.Block
		case <-time.After(chaindb_config.BlockTime):
			return
		}

		candidates = append(candidates, previous_blk)
		previous_uuid = previous_blk.ParentUUID
	}

	for _, c := range candidates {
		blockchain.BlockCache.Store(c.UUID, c)
	}
	if chain.TrySwitchHead(head.UUID) {
		log.Printf(
			"[ Switch] 已切换到拉取自网络的更长链, HEAD.UUID=%d, HEAD.Height=%d\n",
			head.UUID,
			head.Height,
		)
	}
	go func() {
		for _, c := range candidates {
			blockCacheDetached.Delete(c.UUID)
		}
	}()
}
