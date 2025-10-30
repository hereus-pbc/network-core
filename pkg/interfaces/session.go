package interfaces

import "github.com/hereus-pbc/network-core/pkg/types"

type SessionManager interface {
	Create(user User, loginMethod string, pass string, packageName string) (string, error)
	Exists(sessionToken string) (bool, error)
	Get(sessionToken string) (*types.SessionsModal, error)
}

type Session interface {
	GetKernel() Kernel
	GetSession() *types.SessionsModal
	GetUser() User
	GetSessionToken() string
}
