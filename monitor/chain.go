package monitor

import (
	"encoding/json"
	"net/http"
)

func registerChainService() {
	http.HandleFunc("/head", getHeadHandler)
}

// GET /head
func getHeadHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(theChain.Head().UUID)
}
