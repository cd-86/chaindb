package main

import (
	"fmt"
	"time"

	"github.com/shynur/chaindb/discovery"
)

func main() {
	discovery.Start(1 * time.Second)
	for ; ; time.Sleep(1 * time.Second) {
		for _, ip_addr := range *discovery.ActiveNodes.Load() {
			fmt.Println(ip_addr)
		}
		fmt.Println()
	}
}
