package main

import (
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/shynur/chaindb/chaindb_config"
)

func StartUserService() {
	server := http.NewServeMux()

	server.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {

	})

	go func() {
		err := http.ListenAndServe(
			net.JoinHostPort(
				"",
				strconv.Itoa(chaindb_config.TCPPortUserService),
			),
			server,
		)
		log.Fatalln("[UserService] Error:", err)
	}()
}
