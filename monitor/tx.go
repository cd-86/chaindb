package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
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

// /owners
func ownersHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]uint32{})
}

// /owners/{owner_id}/transactions
func owners_TransactionsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.NotFound(w, r)
		return
	}
	ownerID := parts[2]
	var id int
	_, err := fmt.Sscanf(ownerID, "%d", &id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	mtx.RLock()
	defer mtx.RUnlock()
	json.NewEncoder(w).Encode(transactionStore[id])
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
