package monitor

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/shynur/chaindb/discovery"
)

func registerNodeService() {
	http.HandleFunc("/peers", peersHandler)
	http.HandleFunc("/peers/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			deletePeersHostHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}

// DELETE /peers/[host]
func deletePeersHostHandler(w http.ResponseWriter, r *http.Request) {
	host := strings.Split(r.URL.Path, "/")[2]
	discovery.AdministratorSpecifiedNodes.Remove([]string{host})
	w.WriteHeader(http.StatusNoContent)
}

func peersHandler(w http.ResponseWriter, r *http.Request) {
	// GET /peers
	get := func() {
		json.NewEncoder(w).Encode(
			slices.Collect(
				func(yield func(v struct {
					Host                   string
					AdministratorSpecified bool
				}) bool) {
					admin_specified := discovery.AdministratorSpecifiedNodes.List()
					for _, host := range admin_specified {
						yield(struct {
							Host                   string
							AdministratorSpecified bool
						}{
							Host:                   host,
							AdministratorSpecified: true,
						})
					}

					for _, peer := range *discovery.ActiveNodes.Load() {
						if found := slices.Index(admin_specified, peer); found != -1 {
							admin_specified = slices.Delete(
								admin_specified,
								found, found+1,
							)
							continue
						}
						yield(struct {
							Host                   string
							AdministratorSpecified bool
						}{
							Host:                   peer,
							AdministratorSpecified: false,
						})
					}
				},
			),
		)
	}

	// POST /peers
	post := func() {
		var new_nodes []string
		if err := json.NewDecoder(r.Body).Decode(&new_nodes); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		discovery.AdministratorSpecifiedNodes.Add(new_nodes)
		w.WriteHeader(http.StatusCreated)
	}

	switch r.Method {
	case http.MethodGet:
		get()
	case http.MethodPost:
		post()
	default:
		http.NotFound(w, r)
	}
}
