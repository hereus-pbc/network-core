package rpc_net_hereus_sdk_activitypub_actor

import (
	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func ListFollowing(session interfaces.Session) *[]HandleActorResponse {
	var resp []HandleActorResponse
	for _, followingUrl := range session.GetUser().ListFollowing() {
		following, err := session.GetKernel().ActivityPubDB().GetActor(followingUrl, "")
		if err != nil {
			continue
		}
		followingHandle, err := session.GetKernel().ReverseActorUrl(followingUrl, following)
		if err == nil {
			resp = append(resp, ActorToResponse(following, followingHandle))
		}
	}
	if len(resp) == 0 {
		return nil
	}
	return &resp
}
