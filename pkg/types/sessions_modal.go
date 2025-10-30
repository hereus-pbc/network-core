package types

type SessionsModal struct {
	SessionTokenHash string `bson:"session_token_hash"`
	Username         string `bson:"username"`
	LoginMethod      string `bson:"login_method"` // "password", "passkey", "eth", "theprotocols", "token"
	CreatedAt        string `bson:"created_at"`
	LastAccessed     string `bson:"last_accessed"`
	PackageName      string `bson:"package_name"`
	MidKey           string `bson:"mid_key,omitempty"`
}
