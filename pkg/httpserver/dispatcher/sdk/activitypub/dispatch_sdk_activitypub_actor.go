package dispatcher_sdk_activitypub

import (
	"fmt"
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	rpc_net_hereus_sdk_activitypub_actor "github.com/hereus-pbc/network-core/pkg/rpcserver/net/hereus/sdk/activitypub/actor"
)

func handleActivityPubSDKActor(session interfaces.Session, w http.ResponseWriter, r *http.Request, innerFunction string) (out interface{}, err error) {
	switch innerFunction {
	case "follow":
		var req rpc_net_hereus_sdk_activitypub_actor.FollowArguments
		if helpers.TryReadBody(r, w, &req) {
			err = rpc_net_hereus_sdk_activitypub_actor.Follow(session, &req)
		}
	case "unfollow":
		var req rpc_net_hereus_sdk_activitypub_actor.UnfollowArguments
		if helpers.TryReadBody(r, w, &req) {
			err = rpc_net_hereus_sdk_activitypub_actor.Unfollow(session, &req)
		}
	case "get":
		var req rpc_net_hereus_sdk_activitypub_actor.GetArguments
		if helpers.TryReadBody(r, w, &req) {
			return rpc_net_hereus_sdk_activitypub_actor.Get(session, &req)
		}
	case "listFollowers":
		return rpc_net_hereus_sdk_activitypub_actor.ListFollowers(session), nil
	case "listFollowing":
		return rpc_net_hereus_sdk_activitypub_actor.ListFollowing(session), nil
	default:
		return nil, fmt.Errorf("unknown activity endpoint: net.hereus.sdk.activitypub.actor.%s", innerFunction)
	}
	return nil, nil
}
