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
			Name:               "HereUS Network",
			Icon:               "https://static.hereus.net/favicon-nobg.png",
			Description:        "",
			LatestVersion:      kernel.GetSoftwareVersion(),
			LatestBuild:        kernel.GetSoftwareBuild(),
			DefaultPreferences: map[string]interface{}{},
			InitialPermissions: []string{
				"Security.GenerateSessions",
				"Security.AccessHiddenInformation",
			},
		}
		appInfoBytes, err := json.Marshal(appInfo)
		if err != nil {
			http.Error(w, "Failed to generate app info", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(appInfoBytes)
	}
}
