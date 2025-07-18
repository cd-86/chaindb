package main

import (
	"fmt"
	"time"

	"github.com/shynur/chaindb/discovery"
)

func main() {
	discovery.Start(1 * time.Second)
	for ; ; time.Sleep(1 * time.Second) {
		fmt.Println()
	}
}
