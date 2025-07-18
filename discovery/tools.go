package discovery

import (
	"net"
)

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

func isLocalIP(ip_addr net.IP) bool {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			var local_addr net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				local_addr = v.IP
			case *net.IPAddr:
				local_addr = v.IP
			}
			if local_addr != nil && ip_addr.Equal(local_addr) {
				return true
			}
		}
	}

	return false
}
