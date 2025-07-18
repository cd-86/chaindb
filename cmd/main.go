package main

import (
	"time"

	"github.com/shynur/chaindb/discovery"
)

func main() {
	/* for _, iface := range discovery.GetOneMulticastNetworkInterface() {
		fmt.Println(
			iface.MulticastAddrs(),
		)
	} */
	discovery.Register()
	discovery.FindAll(5 * time.Second)
}
