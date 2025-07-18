package main

import (
	"fmt"
	"time"

	"github.com/shynur/chaindb/discovery"
)

func main() {
	prelude()
	discovery.Start(1 * time.Second)
	for ; ; time.Sleep(1 * time.Second) {
		fmt.Println()
	}
}
