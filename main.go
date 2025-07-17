package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/mdns"
)

func s() {
	// Setup our service export
	host, _ := os.Hostname()
	info := []string{"My awesome service"}
	service, _ := mdns.NewMDNSService(host, "_foobar._tcp", "", "", 8000, nil, info)

	// Create the mDNS server, defer shutdown
	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	defer server.Shutdown()
	time.Sleep(10 * time.Second)
}
func r() {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("!!! Got new entry: %v\n", entry)
		}
	}()

	// Start the lookup
	params := mdns.DefaultParams("_foobar._tcp")
	params.Entries = entriesCh
	params.DisableIPv6 = true
	params.Timeout = 5 * time.Second
	mdns.Query(params)
	close(entriesCh)
}
func main() {
	go s()
	go r()
	time.Sleep(5 * time.Second)
}
