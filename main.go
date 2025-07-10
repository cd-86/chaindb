package main

import (
	"fmt"
	"time"

	"github.com/hashicorp/mdns"
)

func publish() {
	service, _ := mdns.NewMDNSService(
		"叫啥名字应该无所谓吧", "shynur.ChainDB.Discovery", "", "", 8000, nil,
		nil,
	)

	// Create the mDNS server, defer shutdown
	mdns.NewServer(&mdns.Config{Zone: service, LogEmptyResponses: true})
	// defer server.Shutdown()
}
func lookup() {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("!!! %v\n", entry.AddrV4)
		}
	}()

	mdns.Lookup("shynur.ChainDB.Discovery", entriesCh)
}

func main() {
	publish()
	lookup()
	time.Sleep(3 * time.Second)
}
