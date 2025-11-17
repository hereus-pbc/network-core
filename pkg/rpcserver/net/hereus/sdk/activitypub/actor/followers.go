package rpc_net_hereus_sdk_activitypub_actor

import (
	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func ListFollowersRpc() *helpers.RpcFunctionHandlerNoArguments {
	return &helpers.RpcFunctionHandlerNoArguments{
		Handler: func(session interfaces.Session) (interface{}, error) {
			return ListFollowers(session), nil
		},
		Permissions: []string{"net.hereus.sdk.permissions.activitypub"},
	}
}

func ListFollowers(session interfaces.Session) *[]HandleActorResponse {
	var resp []HandleActorResponse
	for _, followerUrl := range session.GetUser().ListFollowers() {
		follower, err := session.GetKernel().ActivityPubDB().GetActor(followerUrl, "")
		if err != nil {
			continue
		}
		followerHandle, err := session.GetKernel().ReverseActorUrl(followerUrl, follower)
		if err == nil {
			resp = append(resp, ActorToResponse(follower, followerHandle))
		}
	}
	if len(resp) == 0 {
		return nil
	}
	return &resp
}
