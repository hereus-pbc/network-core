package types

type AppInfo struct {
	Name               string                 `json:"name"`
	Icon               string                 `json:"icon"`
	Description        string                 `json:"description"`
	LatestVersion      string                 `json:"latestVersion"`
	LatestBuild        int                    `json:"latestBuildNumber"`
	DefaultPreferences map[string]interface{} `json:"defaultPreferences"`
	InitialPermissions []string               `json:"initialPermissions"`
}
