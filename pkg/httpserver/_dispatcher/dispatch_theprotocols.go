package http_server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
	rpc_org_theprotocols_network "github.com/hereus-pbc/network-core/pkg/rpcserver/org/theprotocols/network"
	rpc_org_theprotocols_session "github.com/hereus-pbc/network-core/pkg/rpcserver/org/theprotocols/session"
)

func handleTheProtocolsCore(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request, functionName string, endpoint string) {
	switch functionName {
	case "network":
		out, err := rpc_org_theprotocols_network.HandleNetworkInformation(kernel)
		convertRpcResponseToHttpResponse(out, err, w)
	case "session.getUserId":
		kernel.SessionWrapper(w, r, []string{}, func(session interfaces.Session) {
			out, err := rpc_org_theprotocols_session.GetUserId(session)
			convertRpcResponseToHttpResponse(out, err, w)
		})
	default:
		http.Error(w, fmt.Sprintf("Unknown endpoint: %s", endpoint), http.StatusNotFound)
	}
}

func dispatchRpcEndpoints(kernel interfaces.Kernel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Request-Headers", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		endpoint, _ := strings.CutPrefix(r.URL.Path, "/theprotocols/")
		if len(endpoint) == 0 {
			http.Error(w, "Endpoint not specified", http.StatusBadRequest)
			return
		}
		fmt.Printf("(TheProtocols) %s [%s] %s\n", r.RemoteAddr, time.Now().Format("2006-01-02 15:04:05"), endpoint)
		endpointPieces := strings.Split(endpoint, ".")
		if len(endpointPieces) < 3 {
			http.Error(w, "Invalid endpoint format", http.StatusBadRequest)
			return
		}
		domainRoot := fmt.Sprintf("%s.%s", endpointPieces[1], endpointPieces[0])
		// collect all left together
		functionName := strings.Join(endpointPieces[2:], ".")
		switch domainRoot {
		case "theprotocols.org":
			handleTheProtocolsCore(kernel, w, r, functionName, endpoint)
		case "hereus.net":
			handleHereUS(kernel, w, r, functionName, endpoint)
		default:
			http.Error(w, fmt.Sprintf("Unknown domain root: %s", domainRoot), http.StatusNotFound)
		}
	}
}
