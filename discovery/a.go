package discovery

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
	"slices"
	"sync"
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

func FindAll(timeout time.Duration) []net.IP {
	var (
		nodes_list_lock  sync.Mutex
		discovered_nodes []net.IP
	)

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if entry.Instance == UniqueNodeName {
				continue
			}

			var has_been_added atomic.Bool
			for _, addr := range append(entry.AddrIPv4, entry.AddrIPv6...) {
				go func() {
					if err := ping(addr.String()); err != nil {
						return
					}
					if has_been_added.CompareAndSwap(false, true) {
						nodes_list_lock.Lock()
						defer nodes_list_lock.Unlock()
						discovered_nodes = append(discovered_nodes, addr)
						log.Printf("已发现 %s\n", addr)
					}
				}()
			}
		}
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

	nodes_list_lock.Lock()
	defer nodes_list_lock.Unlock()
	return discovered_nodes
}

// INTERVAL: 网络环境中检查各类事项的间隔时间或超时时间.
func Start(interval time.Duration) {
	Register(interval)
	go func() {
		nodes := FindAll(interval)
		ActiveNodes.Store(&nodes)
	}()
}

func Register(interval_checking_network time.Duration) {
	ifcs := []net.Interface{}
	ifaces_is_equal := func(ifs1, ifs2 []net.Interface) bool {
		if len(ifs1) != len(ifs2) {
			return false
		}
		for _, iface1 := range ifs1 {
			if !slices.ContainsFunc(
				ifs2,
				func(iface2 net.Interface) bool {
					if iface1.Index != iface2.Index {
						return false
					}
					if iface1.MTU != iface2.MTU {
						return false
					}
					if iface1.Name != iface2.Name {
						return false
					}
					if iface1.HardwareAddr.String() != iface2.HardwareAddr.String() {
						return false
					}
					if iface1.Flags != iface2.Flags {
						return false
					}
					return true
				},
			) {
				return false
			}
		}
		return true
	}

	var server *zeroconf.Server
	go func() {
		for ; ; time.Sleep(interval_checking_network) {
			new_ifcs := getMulticastNetworkInterfaces()
			if ifaces_is_equal(ifcs, new_ifcs) {
				continue
			}

			if server != nil {
				server.Shutdown()
			}
			server, _ = zeroconf.Register(
				UniqueNodeName,
				"_shynur-chaindb._tcp", "local.", chaindb.DNSSDPort,
				nil, new_ifcs,
			)

			ifcs = new_ifcs
		}
	}()
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
