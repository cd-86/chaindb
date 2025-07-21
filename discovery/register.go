package discovery

import (
	"net"
	"slices"
	"time"

	"github.com/grandcat/zeroconf"
	chaindb_config "github.com/shynur/chaindb/config"
)

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
				"_shynur-chaindb._tcp", "local.", chaindb_config.DNSSDPort,
				nil, new_ifcs,
			)

			ifcs = new_ifcs
		}
	}()
}
