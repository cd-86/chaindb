package validator

import (
	"context"
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
	_ = blockchain.BlockCache

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

	server := grpc.NewServer()
	block_cdn.RegisterBlockDeliveryServer(
		server,
		&blockPCDNServer{},
	)
	if err := server.Serve(listener); err != nil {
		panic(err)
	}
}

type blockPCDNServer struct {
	block_cdn.UnimplementedBlockDeliveryServer
}

func (s blockPCDNServer) GetBlock(
	ctx context.Context,
	arg *block_cdn.BlockUUID,
) (*block_cdn.BlockGob, error) {
	blk, ok := blockchain.BlockCache.Load(arg.GetUUID())

	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"block not found: %d", arg.UUID,
		)
	}

	return &block_cdn.BlockGob{
		Buffer: blk.(blockchain.Block).ToGob(),
	}, nil
}
