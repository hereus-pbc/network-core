package dispatcher_theprotocols

import (
	"fmt"
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	rpc_org_theprotocols_network "github.com/hereus-pbc/network-core/pkg/rpcserver/org/theprotocols/network"
	rpc_org_theprotocols_session "github.com/hereus-pbc/network-core/pkg/rpcserver/org/theprotocols/session"
)

func HandleTheProtocolsCore(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request, functionName string, endpoint string) {
	switch functionName {
	case "network":
		out, err := rpc_org_theprotocols_network.HandleNetworkInformation(kernel)
		helpers.ConvertRpcResponseToHttpResponse(out, err, w)
	case "session.getUserId":
		kernel.SessionWrapper(w, r, []string{}, func(session interfaces.Session) {
			out, err := rpc_org_theprotocols_session.GetUserId(session)
			helpers.ConvertRpcResponseToHttpResponse(out, err, w)
		})
	default:
		http.Error(w, fmt.Sprintf("Unknown endpoint: %s", endpoint), http.StatusNotFound)
	}
}
