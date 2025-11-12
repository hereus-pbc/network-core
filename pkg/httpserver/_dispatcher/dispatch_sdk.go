package http_server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hereus-pbc/network-core/pkg/httpserver/dispatcher/sdk/activitypub"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func handleHereUS(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request, functionName string, endpoint string) {
	if strings.HasPrefix(functionName, "sdk.activitypub.") {
		dispatcher_sdk_activitypub.HandleActivityPubSDK(kernel, w, r, functionName, endpoint)
		return
	}
	http.Error(w, fmt.Sprintf("Unknown endpoint: %s", endpoint), http.StatusNotFound)
}
