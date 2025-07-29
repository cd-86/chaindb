package monitor

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/shynur/chaindb/blockchain"
)

func registerBlockService(server *http.ServeMux) {
	server.HandleFunc("/blocks/", getBlocksUUIDHandler)
}

// GET /blocks/[uuid]
func getBlocksUUIDHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.NotFound(w, r)
		return
	}
	uuid, err := strconv.Atoi(parts[2])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	block, ok := blockchain.BlockCache.Load(uint32(uuid))
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Cache-Control", "public, immutable, max-age=31536000")
	json.NewEncoder(w).Encode(block)
}
