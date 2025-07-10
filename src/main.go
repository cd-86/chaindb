package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/mdns"
)

func publish() {
	// Setup our service export
	host, _ := os.Hostname()
	fmt.Println("host 名字:", host)

	service, _ := mdns.NewMDNSService(
		host, "_foobar._tcp", "", "", 8000, nil,
		[]string{"My awesome service"},
	)

	// Create the mDNS server, defer shutdown
	mdns.NewServer(&mdns.Config{Zone: service})
	// defer server.Shutdown()
	time.Sleep(100 * time.Second)
}
func lookup() {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %v\n", entry)
		}
	}()

	// Start the lookup
	qparams := mdns.DefaultParams("_foobar._tcp")
	qparams.Entries = entriesCh
	qparams.Timeout = 5 * time.Second
	qparams.DisableIPv6 = true
	mdns.Query(qparams)
	time.Sleep(100 * time.Second)
}

func main() {
	go publish()
	time.Sleep(time.Second)
	go lookup()
	time.Sleep(5 * time.Second)
}
