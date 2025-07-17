package discovery

import (
	"context"
	"fmt"
	"net"

	"github.com/pion/mdns/v2"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func Query() {
	addr_ipv4, err := net.ResolveUDPAddr("udp4", mdns.DefaultAddressIPv4)
	if err != nil {
		panic(err)
	}
	socket_ipv4, err := net.ListenUDP("udp4", addr_ipv4)
	if err != nil {
		panic(err)
	}

	addr_ipv6, err := net.ResolveUDPAddr("udp6", mdns.DefaultAddressIPv6)
	if err != nil {
		panic(err)
	}
	socket_ipv6, err := net.ListenUDP("udp6", addr_ipv6)
	if err != nil {
		panic(err)
	}

	mdns_node, err := mdns.Server(
		ipv4.NewPacketConn(socket_ipv4),
		ipv6.NewPacketConn(socket_ipv6),
		&mdns.Config{LocalNames: []string{"shynur-chaindb.local"}},
	)
	if err != nil {
		panic(err)
	}

	answer, src, err := mdns_node.QueryAddr(context.TODO(), "shynur-chaindb.local")

	fmt.Println(answer)
	fmt.Println(src)
	fmt.Println(err)
}
