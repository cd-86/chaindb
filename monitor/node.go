package monitor

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/shynur/chaindb/discovery"
)

func registerNodeService() {
	http.HandleFunc("/peers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			peerListHandler(w, r)
		case http.MethodPost:
			peerCreateHandler(w, r)
		default:
			http.NotFound(w, r)
		}
	})
	http.HandleFunc("/peers/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			peerDeleteHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}

// GET /peers
func peersHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(
		slices.Collect(
			func(yield func(v struct {
				Host                   string
				AdministratorSpecified bool
			}) bool) {
				admin_specified := discovery.AdministratorSpecifiedNodes.List()

				for peer := range *discovery.ActiveNodes.Load() {

				}
			},
		),
	)
}

// POST /peers
func peerCreateHandler(w http.ResponseWriter, r *http.Request) {
	var p Peer
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil || p.Host == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	mtx.Lock()
	peerStore[p.Host] = p
	mtx.Unlock()
	w.WriteHeader(http.StatusCreated)
}

// DELETE /peers/{host}
func peerDeleteHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.NotFound(w, r)
		return
	}
	host := parts[2]
	mtx.Lock()
	delete(peerStore, host)
	mtx.Unlock()
	w.WriteHeader(http.StatusNoContent)
}
