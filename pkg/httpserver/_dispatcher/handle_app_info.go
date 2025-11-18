package http_server

import (
	"encoding/json"
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

func handleAppInfo(kernel interfaces.Kernel) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		appInfo := types.AppInfo{
			Name:               kernel.GetName(),
			Icon:               kernel.GetIcon(),
			Description:        "",
			LatestVersion:      kernel.GetSoftwareVersion(),
			LatestBuild:        kernel.GetSoftwareBuild(),
			DefaultPreferences: map[string]interface{}{},
			InitialPermissions: []string{
				"org.theprotocols.session.permission.security.accessHiddenInformation",
				"org.theprotocols.session.permission.security.generateSessions",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(appInfo); err != nil {
			http.Error(w, "Failed to encode app info", http.StatusInternalServerError)
			return
		}
	}
}
