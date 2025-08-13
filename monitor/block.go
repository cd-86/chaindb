package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/validator"
)

func registerBlockService(server *http.ServeMux) {
	server.HandleFunc("/blocks", blocksHandler)
	server.HandleFunc("/blocks/", getBlocksUUIDHandler)
}

func blocksHandler(w http.ResponseWriter, r *http.Request) {
	// POST /blocks
	// JSON: [block1, ...]
	if r.Method != http.MethodPost {
		http.Error(w, "只支持 POST", http.StatusMethodNotAllowed)
		return
	}

	var blocks []blockchain.Block

	err := json.NewDecoder(r.Body).Decode(&blocks)
	if err != nil {
		http.Error(w, fmt.Sprintf("无效的 JSON (error: %s)", err), http.StatusBadRequest)
		return
	}

	for _, blk := range blocks {
		_, exists := blockchain.BlockCache.Load(blk.UUID)
		if exists {
			continue
		}

		validator.BlockCacheDetached.Store(
			blk.UUID,
			blk,
		)
	}
}

// GET /blocks/[uuid]
func getBlocksUUIDHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.NotFound(w, r)
		return
	}
	uuid, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	block, ok := validator.BlockCacheDetached.Load(uint64(uuid))
	if !ok {
		block, ok = blockchain.BlockCache.Load(uint64(uuid))
		if !ok {
			http.NotFound(w, r)
			return
		}
	}

	w.Header().Set("Cache-Control", "public, immutable, max-age=31536000")
	json.NewEncoder(w).Encode(block)
}
