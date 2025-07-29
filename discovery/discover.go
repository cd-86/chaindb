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

func FindAllActive(timeout time.Duration) (hosts []string) {
	var (
		hosts_          []string
		hosts_list_lock sync.Mutex
	)

	for _, admin_specified := range AdministratorSpecifiedNodes.List() {
		go func() {
			if isLocalIP(net.ParseIP(admin_specified)) {
				log.Fatalf(
					"用户指定的节点 <%s> 是本机 IP, 本该在 HTTP 端就被过滤掉才对啊.\n",
					admin_specified,
				)
			}
			if err := ping(admin_specified); err != nil {
				return
			}
			hosts_list_lock.Lock()
			defer hosts_list_lock.Unlock()
			hosts_ = append(hosts_, admin_specified)
			log.Printf("[Monitor] 用户指定节点 <%s> 是可达的\n", admin_specified)
		}()
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if entry.Instance == chaindb_config.MinerAddress {
				continue
			}

			var has_been_added atomic.Bool
			for _, addr := range slices.Concat(
				entry.AddrIPv4,
				// entry.AddrIPv6,
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
	err = resolver.Browse(ctx, "_shynurChainDB._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err)
	}

	<-ctx.Done()

	hosts_list_lock.Lock()
	defer hosts_list_lock.Unlock()
	return slices.Clone(hosts_)
}
