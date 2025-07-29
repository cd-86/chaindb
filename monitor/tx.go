package monitor

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/shynur/chaindb/blockchain"
)

func registerTxService() {
	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postTransactionsHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	http.HandleFunc("/owners", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/owners" {
			getOwnersHandler(w, r)
			return
		}
		if parts := strings.Split(r.URL.Path, "/"); len(parts) == 4 &&
			parts[1] == "owners" &&
			parts[3] == "transactions" {
			getOwnersOwnerIDTransactionsHandler(w, r)
			return
		}
		http.NotFound(w, r)
	})
}

// GET /owners
func getOwnersHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(
		theChain.ListOwners(),
	)
}

// GET /owners/{owner_id}/transactions
func getOwnersOwnerIDTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	owner, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[2])
	if err != nil {
		http.NotFound(w, r)
		return
	}

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
