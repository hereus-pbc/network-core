package interfaces

import (
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/types"
)

type Kernel interface {
	PushIncomingActivity(activity types.ActivityStream)
	PushOutgoingActivity(activity types.ActivityStream)
	SendOutgoingActivity(activity types.ActivityStream)
	ProcessActivityRoutine()
	SendOutgoingActivityRoutine(user User, sessionToken string, session *types.SessionsModal)
	ReadConfig(key string, result interface{}) error
	ReadConfigString(key string) (string, error)
	GetDomain() string
	SessionManager() SessionManager
	UserManager() UserManager
	BlobManager() BlobManager
	ActivityPubDB() ActivityPubDB
	GetSoftwareVersion() string
	GetSoftwareBuild() int
	GetSoftwareVersionChannel() string
	FetchSession(r *http.Request, fetchList []string) (Session, error)
	SessionWrapper(w http.ResponseWriter, r *http.Request, fetchList []string, handler func(singleFetch Session))
	DoWebfingerRequest(domain string, resource string) (types.WebFingerResponse, error)
	GetActorUrlByHandler(actorHandler string) (string, string, error)
	ReverseActorUrl(actorUrl string, actor types.Actor) (string, error)
}
