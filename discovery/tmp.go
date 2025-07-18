package discovery

import "net"

func GetOneMulticastNetworkInterface() (interfaces []net.Interface) {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, ifi := range ifaces {
		if ifi.Flags&net.FlagUp == 0 {
			continue
		}
		if ifi.Flags&net.FlagMulticast == 0 {
			continue
		}
		if ifi.Flags&net.FlagLoopback > 0 {
			continue
		}
		interfaces = append(interfaces, ifi)
	}
	return interfaces
}
