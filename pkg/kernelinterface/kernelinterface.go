package kernelinterface

import (
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/types"
)

type KernelInterface interface {
	PushIncomingActivity(activity types.ActivityStream)
	PushOutgoingActivity(activity types.ActivityStream)
	SendOutgoingActivity(activity types.ActivityStream)
	ProcessActivityRoutine()
	SendOutgoingActivityRoutine(user *UserInterface, sessionToken string, session *types.SessionsModal)
	ReadConfig(key string, result interface{}) error
	ReadConfigString(key string) (string, error)
	GetDomain() string
	SessionManager() *SessionManagerInterface
	UserManager() *UserManagerInterface
	BlobManager() *BlobManagerInterface
	ActivityPubDB() *ActivityPubDBInterface
	GetSoftwareVersion() string
	GetSoftwareBuild() int
	GetSoftwareVersionChannel() string
	SingleFetcher(r *http.Request, fetchList []string) (*SingleFetcherResultInterface, error)
	SingleFetcherWrapper(w http.ResponseWriter, r *http.Request, fetchList []string, handler func(singleFetch *SingleFetcherResultInterface))
	DoWebfingerRequest(domain string, resource string) (types.WebFingerResponse, error)
	GetActorUrlByHandler(actorHandler string) (string, string, error)
	ReverseActorUrl(actorUrl string, actor types.Actor) (string, error)
}

type UserInterface interface{}

type SessionManagerInterface struct{}

type UserManagerInterface struct{}

type BlobManagerInterface struct{}

type ActivityPubDBInterface struct{}

type SingleFetcherResultInterface struct{}
