package monitor

import (
	"encoding/json"
	"net/http"
)

func registerChainService() {
	http.HandleFunc("/head", headHandler)
}

// GET /head
func headHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(theChain.Head().UUID)
}
