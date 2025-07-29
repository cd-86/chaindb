package monitor

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/shynur/chaindb/blockchain"
)

func registerBlockService() {
	http.HandleFunc("/blocks/", getBlocksUUIDHandler)
}

// GET /blocks/[uuid]
func getBlocksUUIDHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.NotFound(w, r)
		return
	}
	uuid := parts[2]

	block, ok := blockchain.BlockCache.Load(uuid)
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Cache-Control", "public, immutable, max-age=31536000")
	json.NewEncoder(w).Encode(block)
}
