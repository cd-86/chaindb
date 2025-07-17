package main

import (
	"time"

	"github.com/shynur/chaindb/chaindb/discovery"
)

func main() {
	discovery.Register()
	discovery.FindAll(10 * time.Second)
}
