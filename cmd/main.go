package main

import (
	"fmt"

	"github.com/shynur/chaindb/discovery"
)

func main() {
	for _, iface := range discovery.GetOneMulticastNetworkInterface() {
		fmt.Println(
			iface,
		)
	}
}
