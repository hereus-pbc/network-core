package interfaces

type App interface {
	CheckPermission(permission string) bool
	GetData(name string) ([]byte, error)
	SaveData(name string, data []byte) error
	GetPreferencesLastUpdate() int64
	GetPreferences() (map[string]interface{}, error)
	SavePreferences(preferences map[string]interface{}) error
}
