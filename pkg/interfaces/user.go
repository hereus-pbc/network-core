package interfaces

import "github.com/hereus-pbc/network-core/pkg/types"

type User interface {
	GetUsername() string
	GetDescription() string
	GetFirstName() string
	GetMiddleName() string
	GetLastName() string
	GetNameSuffix() string
	GetFullName() string
	GetRsaPublicKey() string
	DeriveMiddleKey(encryptionType string, pass string, newPass string) (string, error)
	DecryptWithSession(sessionToken string, session *types.SessionsModal, encrypted string) ([]byte, error)
	GetRsaPrivateKeyPemWithSession(sessionToken string, session *types.SessionsModal) (string, error)
	GetRootActivityResourceUrl(name string) string
	GetActorUrl() string
	ListFollowers() []string
	ListFollowing() []string
	App(packageName string) (App, error)
	CheckPassword(password string) bool
	IsMultiFactorEnabled() bool
	CheckTOTP(code string) bool
	DoAllowEthereumLogin() bool
	AnyPasskeys() bool
	ToUserId() *types.UserId
	BuildActorObject() (types.Actor, error)
}
