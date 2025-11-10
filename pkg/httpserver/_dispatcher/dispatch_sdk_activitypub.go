package http_server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func handleActivityPubSDK(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request, functionName string, endpoint string) {
	kernel.SessionWrapper(w, r, []string{}, func(session interfaces.Session) {
		innerTop := strings.Split(strings.TrimPrefix(functionName, "sdk.activitypub."), ".")[0]
		innerFunction := strings.TrimPrefix(functionName, "sdk.activitypub."+innerTop+".")
		var out interface{}
		var err error
		switch innerTop {
		case "activity":
			out, err = handleActivityPubSDKActivity(session, w, r, innerFunction)
		case "actor":
			out, err = handleActivityPubSDKActor(session, w, r, innerFunction)
		default:
			http.Error(w, fmt.Sprintf("Unknown endpoint: %s", endpoint), http.StatusNotFound)
			return
		}
		convertRpcResponseToHttpResponse(out, err, w)
	})
}
