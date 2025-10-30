package http_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

func handleWebFinger(kernel interfaces.Kernel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resource := r.URL.Query().Get("resource")
		if resource == "" {
			http.Error(w, "Resource parameter is required", http.StatusBadRequest)
			return
		}
		if len(resource) > 5 && resource[:5] == "acct:" {
			resource = resource[5:]
		}
		domain, err := kernel.ReadConfigString("domain")
		if err != nil {
			http.Error(w, "Failed to read domain from config", http.StatusInternalServerError)
			return
		}
		if len(strings.Split(resource, "@")) != 2 {
			http.Error(w, "Invalid resource format", http.StatusBadRequest)
			return
		}
		if !strings.HasSuffix(resource, "@"+domain) {
			http.Error(w, "Resource does not belong to this domain", http.StatusBadRequest)
			return
		}
		username := strings.Split(resource, "@")[0]
		user, err := kernel.UserManager().Get(username)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		webFingerResponse := types.WebFingerResponse{
			Subject: "acct:" + user.GetUsername() + "@" + domain,
			Aliases: []string{
				fmt.Sprintf("https://%s/~%s", domain, user.GetUsername()),
				fmt.Sprintf("https://%s/activitypub/accounts/%s", domain, user.GetUsername()),
			},
			Links: []types.WebFingerLinks{
				{
					Rel:  "self",
					Type: "application/activity+json",
					Href: fmt.Sprintf("https://%s/activitypub/accounts/%s", domain, user.GetUsername()),
				},
				{
					Rel:  "http://webfinger.net/rel/profile-page",
					Type: "text/html",
					Href: fmt.Sprintf("https://%s/~%s", domain, user.GetUsername()),
				},
			},
		}
		w.Header().Set("Content-Type", "application/jrd+json")
		err = json.NewEncoder(w).Encode(webFingerResponse)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
