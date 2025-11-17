package rpc_net_hereus_sdk_activitypub_actor

import (
	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

func GetRpc() *helpers.RpcFunctionHandlerWithArguments {
	return &helpers.RpcFunctionHandlerWithArguments{
		ReqFactory: func() interface{} {
			return &GetArguments{}
		},
		Handler: func(session interfaces.Session, req interface{}) (interface{}, error) {
			return Get(session, req.(*GetArguments))
		},
		Permissions: []string{"net.hereus.sdk.permissions.activitypub"},
	}
}

type GetArguments struct {
	Handle string `json:"handle"` // URL of the actor to handle
}

type HandleActorResponse struct {
	Handle         string `json:"handle"`
	Url            string `json:"url"`
	DisplayName    string `json:"displayName"`
	Biography      string `json:"biography"`
	ProfilePicture string `json:"profilePicture"`
	BannerPicture  string `json:"bannerPicture"`
}

func ActorToResponse(actor types.Actor, handle string) HandleActorResponse {
	return HandleActorResponse{
		Handle:         handle,
		Url:            actor.Id,
		DisplayName:    actor.Name,
		Biography:      actor.Summary,
		ProfilePicture: actor.Icon.Url,
		BannerPicture:  actor.Image.Url,
	}
}

func Get(session interfaces.Session, req *GetArguments) (*HandleActorResponse, error) {
	_, url, err := session.GetKernel().GetActorUrlByHandler(req.Handle)
	if err != nil {
		return nil, err
	}
	actor, err := session.GetKernel().ActivityPubDB().GetActor(url, "")
	if err != nil {
		return nil, err
	}
	resp := ActorToResponse(actor, req.Handle)
	return &resp, nil
}
