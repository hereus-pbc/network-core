package http_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func tryReadBody(r *http.Request, w http.ResponseWriter, req interface{}) bool {
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return false
	}
	return true
}

func handleHereUS(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request, functionName string, endpoint string) {
	if strings.HasPrefix(functionName, "sdk.activitypub.") {
		handleActivityPubSDK(kernel, w, r, functionName, endpoint)
		return
	}
	http.Error(w, fmt.Sprintf("Unknown endpoint: %s", endpoint), http.StatusNotFound)
}
