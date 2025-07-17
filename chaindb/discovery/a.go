package discovery

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/shynur/chaindb/chaindb"
)

//var ActiveNodes []net.IP

func FindAll(wait_time time.Duration) {
	//discovered_nodes := []net.IP{}
	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Println(entry)
		}
		log.Println("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		func() time.Duration {
			if wait_time == 0 {
				return chaindb.BlockTime
			} else {
				return wait_time
			}
		}(),
	)
	defer cancel()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err)
	}
	err = resolver.Browse(ctx, "_shynur-chaindb._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err)
	}

	<-ctx.Done()
}

func Register() {
	_, err := zeroconf.Register(
		fmt.Sprintf("ChainDB-No%d", time.Now().UnixMilli()),
		"_shynur-chaindb._tcp", "local.", chaindb.DNSSDPort,
		nil, nil,
	)
	if err != nil {
		panic(err)
	}
}
