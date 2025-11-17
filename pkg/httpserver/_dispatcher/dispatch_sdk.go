package http_server

import (
	"fmt"
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	rpc_net_hereus_sdk_activitypub_activity "github.com/hereus-pbc/network-core/pkg/rpcserver/net/hereus/sdk/activitypub/activity"
	rpc_net_hereus_sdk_activitypub_actor "github.com/hereus-pbc/network-core/pkg/rpcserver/net/hereus/sdk/activitypub/actor"
)

func handleHereUS(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request, functionName string) {
	endpoints := map[string]helpers.RpcFunctionHandler{
		"sdk.activitypub.activity.create":       rpc_net_hereus_sdk_activitypub_activity.CreateRpc(),
		"sdk.activitypub.activity.announce":     rpc_net_hereus_sdk_activitypub_activity.AnnounceRpc(),
		"sdk.activitypub.activity.undoAnnounce": rpc_net_hereus_sdk_activitypub_activity.UndoAnnounceRpc(),
		"sdk.activitypub.activity.delete":       rpc_net_hereus_sdk_activitypub_activity.DeleteRpc(),
		"sdk.activitypub.activity.get":          rpc_net_hereus_sdk_activitypub_activity.GetRpc(),
		"sdk.activitypub.activity.list":         rpc_net_hereus_sdk_activitypub_activity.ListRpc(),
		"sdk.activitypub.activity.like":         rpc_net_hereus_sdk_activitypub_activity.LikeRpc(),
		"sdk.activitypub.activity.unlike":       rpc_net_hereus_sdk_activitypub_activity.UnlikeRpc(),
		"sdk.activitypub.activity.edit":         rpc_net_hereus_sdk_activitypub_activity.EditRpc(),
		"sdk.activitypub.actor.follow":          rpc_net_hereus_sdk_activitypub_actor.FollowRpc(),
		"sdk.activitypub.actor.unfollow":        rpc_net_hereus_sdk_activitypub_actor.UnfollowRpc(),
		"sdk.activitypub.actor.get":             rpc_net_hereus_sdk_activitypub_actor.GetRpc(),
		"sdk.activitypub.actor.listFollowers":   rpc_net_hereus_sdk_activitypub_actor.ListFollowersRpc(),
		"sdk.activitypub.actor.listFollowing":   rpc_net_hereus_sdk_activitypub_actor.ListFollowingRpc(),
	}
	if handler, exists := endpoints[functionName]; exists {
		handler.Handle(kernel, r, w)
		return
	}
	http.Error(w, fmt.Sprintf("Unknown endpoint: net.hereus.%s", functionName), http.StatusNotFound)
}
