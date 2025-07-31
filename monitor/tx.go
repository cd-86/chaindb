package monitor

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/shynur/chaindb/blockchain"
)

func registerTxService(server *http.ServeMux) {
	server.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			postTransactionsHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	server.HandleFunc("/owners", getOwnersHandler)
	server.HandleFunc("/owners/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")

		if len(parts) == 4 &&
			parts[1] == "owners" &&
			parts[3] == "transactions" {
			getOwnersOwnerIDTransactionsHandler(w, r)
			return
		}

		log.Printf("[Monitor] 无效的路径: %s\n", r.URL.Path)
		http.NotFound(w, r)
	})
}

// GET /owners
func getOwnersHandler(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(
		theChain.ListOwners(),
	)
}

// GET /owners/{owner_id}/transactions
func getOwnersOwnerIDTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	owner, err := strconv.ParseUint(
		strings.Split(r.URL.Path, "/")[2],
		10, 32,
	)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	log.Printf("[Monitor] 被请求查询的交易记录 OwnerID=%d\n", owner)
	json.NewEncoder(w).Encode(
		theChain.ListTxsOwnedBy(uint32(owner)),
	)
}

// POST /transactions
func postTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	var tx blockchain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "无法解析的 transaction", http.StatusBadRequest)
		return
	}

	theTxPool.Add(tx)
	w.WriteHeader(http.StatusCreated)
}
