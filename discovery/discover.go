package discovery

import (
	"context"
	"log"
	"net"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/shynur/chaindb/chaindb_config"
)

func FindAll(timeout time.Duration) (hosts []string) {
	var (
		hosts_          []string
		hosts_list_lock sync.Mutex
	)

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if entry.Instance == chaindb_config.MinerAddress {
				continue
			}

			var has_been_added atomic.Bool
			for _, addr := range append(
				entry.AddrIPv4,
				[]net.IP{}...,
			// entry.AddrIPv6...,  // 路由器可能不支持 IPv6 吧, 主要是我不确定 'IPv6:port' 咋正确书写.
			) {
				if isLocalIP(addr) {
					continue
				}
				go func() {
					if has_been_added.Load() {
						return
					}
					if err := ping(addr.String()); err != nil {
						return
					}
					if has_been_added.CompareAndSwap(false, true) {
						hosts_list_lock.Lock()
						defer hosts_list_lock.Unlock()
						hosts_ = append(hosts_, addr.String())
						log.Printf("[ DNS-SD] 已发现 %s\n", addr)
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

	hosts_list_lock.Lock()
	defer hosts_list_lock.Unlock()
	return slices.Clone(hosts_)
}
