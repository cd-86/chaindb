package validator

import (
	"context"
	"log"
	"math"
	"net"
	"strconv"

	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/chaindb_config"
	"github.com/shynur/chaindb/validator/block_cdn"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func StartBlockDeliveryServer() {
	log.Printf("[  PCDN ] 正在启动区块 PCDN 服务器...")

	listener, err := net.Listen(
		"tcp",
		net.JoinHostPort(
			"",
			strconv.Itoa(chaindb_config.TCPPortBlockPCDN),
		),
	)
	if err != nil {
		panic(err)
	}
	log.Printf("[  PCDN ] 区块 PCDN 服务器开始监听端口: %d", chaindb_config.TCPPortBlockPCDN)

	server := grpc.NewServer(
		// 放开 gRPC 的所有限制:
		// grpc.InitialConnWindowSize(math.MaxInt32),
		// grpc.InitialWindowSize(math.MaxInt32),
		grpc.MaxHeaderListSize(math.MaxUint32),
		grpc.MaxRecvMsgSize(math.MaxInt32),
		// grpc.StaticConnWindowSize(math.MaxInt32),
		// grpc.StaticStreamWindowSize(math.MaxInt32),
	)
	block_cdn.RegisterBlockDeliveryServer(
		server,
		&blockPCDNServer{},
	)
	go server.Serve(listener)

	log.Printf("[  PCDN ] 区块 PCDN 服务器正常运行")
}

type blockPCDNServer struct {
	block_cdn.UnimplementedBlockDeliveryServer
}

func (s blockPCDNServer) GetBlock(
	ctx context.Context,
	arg *block_cdn.BlockUUID,
) (*block_cdn.BlockGob, error) {
	blk, ok := blockCacheDetached.Load(arg.GetUUID())
	if !ok {
		blk, ok = blockchain.BlockCache.Load(arg.GetUUID())
		if !ok {
			return nil, status.Errorf(
				codes.NotFound,
				"block not found: %d", arg.UUID,
			)
		}
	}

	return &block_cdn.BlockGob{
		Buffer: blk.(blockchain.Block).ToGob(),
	}, nil
}
