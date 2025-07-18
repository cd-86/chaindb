package main

import (
	"time"

	"github.com/shynur/chaindb/discovery"
)

func main() {
	discovery.Register(10 * time.Second)
	discovery.FindAll(5 * time.Second)
}
