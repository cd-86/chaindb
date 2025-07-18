package discovery

import (
	"context"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/grandcat/zeroconf"
)

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
				if isLocalIP(addr) {
					continue
				}
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
