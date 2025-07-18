package discovery

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/shynur/chaindb"
)

var ActiveNodes = func() *atomic.Pointer[[]net.IP] {
	var p atomic.Pointer[[]net.IP]
	p.Store(&[]net.IP{})
	return &p
}() // 不包括自己.

var UniqueNodeName = fmt.Sprintf("ChainDB-No%d", time.Now().UnixMilli())

func FindAll(timeout time.Duration) {
	discovered_nodes := []net.IP{}
	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if entry.Instance == UniqueNodeName {
				continue
			}
			log.Printf("%+v\n", entry)
			for _, addr := range append(entry.AddrIPv4, entry.AddrIPv6...) {
				//net.D
				discovered_nodes = append(discovered_nodes, addr)
			}
		}
		log.Println("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err)
	}
	err = resolver.Browse(ctx, "_shynur-chaindb._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err)
	}

	<-ctx.Done()
	ActiveNodes.Store(&discovered_nodes)
}

func Register() {
	_, err := zeroconf.Register(
		UniqueNodeName,
		"_shynur-chaindb._tcp", "local.", chaindb.DNSSDPort,
		nil, getMulticastNetworkInterfaces(),
	)
	if err != nil {
		panic(err)
	}
}

// (改编自 <https://github.com/grandcat/zeroconf/blob/e4f60f8407b11e9ba16f4c4c5ad24226dd4e8519/connection.go#L101>.)
// 可能有多个, 包括接入的虚拟局域网.
func getMulticastNetworkInterfaces() []net.Interface {
	all_ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	ifaces := []net.Interface{}
	for _, ifi := range all_ifaces {
		if ifi.Flags&net.FlagUp == 0 {
			continue
		}
		if ifi.Flags&net.FlagMulticast == 0 {
			continue
		}
		if ifi.Flags&net.FlagLoopback != 0 {
			continue
		}
		ifaces = append(ifaces, ifi)
	}

	return ifaces
}

func ping(ip string) (err error) {
	switch os := runtime.GOOS; os {
	case "windows":
		err = exec.Command("ping", "-n", "1", ip).Run()
	case "linux":
		err = exec.Command("ping", "-c", "1", ip).Run()
	default:
		panic(
			fmt.Sprintf(
				"平台 OS (%s) 上的 `ping' 暂时没有得到 ChainDB 的支持",
				os,
			),
		)
	}

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return
		} else {
			panic("`exec ping' 失败")
		}
	}
	return
}
