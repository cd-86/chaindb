package monitor

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func registerTxService() {
	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			transactionCreateHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	http.HandleFunc("/owners", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/owners" {
			if strings.HasPrefix(r.URL.Path, "/owners/") && strings.HasSuffix(r.URL.Path, "/transactions") {
				transactionListHandler(w, r)
				return
			}
			http.NotFound(w, r)
			return
		}
		ownerHandler(w, r)
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
func transactionCreateHandler(w http.ResponseWriter, r *http.Request) {
	var t Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	mtx.Lock()
	transactionStore[t.OwnerID] = append(transactionStore[t.OwnerID], t)
	mtx.Unlock()
	w.WriteHeader(http.StatusCreated)
}
