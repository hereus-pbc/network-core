package http_server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	rpc_org_theprotocols "github.com/hereus-pbc/network-core/pkg/rpcserver/org/theprotocols"
	rpc_org_theprotocols_session "github.com/hereus-pbc/network-core/pkg/rpcserver/org/theprotocols/session"
)

func (s *HttpServer) dispatchRpcEndpoints(kernel interfaces.Kernel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Request-Headers", "*")
		if r.Method == http.MethodOptions || r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/theprotocols/") {
			http.Error(w, "Invalid RPC endpoint", http.StatusBadRequest)
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
		domainHandlers := map[string]func(interfaces.Kernel, http.ResponseWriter, *http.Request, string){
			"theprotocols.org": func(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request, functionName string) {
				if helpers.AutoHandleRpcFunction(&helpers.RpcFunctionMapping{
					"network":           rpc_org_theprotocols.NetworkRpc(),
					"session.getUserId": rpc_org_theprotocols_session.GetUserIdRpc(),
				}, kernel, r, w, functionName) != nil {
					http.Error(w, fmt.Sprintf("Unknown endpoint: org.theprotocols.%s", functionName), http.StatusNotFound)
				}
			},
			"hereus.net": handleHereUS,
		}
		if handler, exists := domainHandlers[domainRoot]; exists {
			handler(kernel, w, r, strings.Join(endpointPieces[2:], "."))
			return
		}
		http.Error(w, fmt.Sprintf("Unknown domain root: %s", domainRoot), http.StatusNotFound)
	}
}
