package http_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hereus-pbc/network-core/pkg/httpserver/sdk_ap_handlers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	rpc_net_hereus_sdk_activitypub_activity "github.com/hereus-pbc/network-core/pkg/rpcserver/net/hereus/sdk/activitypub/activity"
	rpc_net_hereus_sdk_activitypub_actor "github.com/hereus-pbc/network-core/pkg/rpcserver/net/hereus/sdk/activitypub/actor"
)

func tryReadBody(r *http.Request, w http.ResponseWriter, req interface{}) bool {
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return false
	}
	return true
}

func handleHereUS(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request, functionName string, endpoint string) {
	if strings.HasPrefix(functionName, "sdk.activitypub.") {
		handleActivityPubSDK(kernel, w, r, functionName, endpoint)
		return
	}
	http.Error(w, fmt.Sprintf("Unknown endpoint: %s", endpoint), http.StatusNotFound)
}

func handleActivityPubSDK(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request, functionName string, endpoint string) {
	kernel.SessionWrapper(w, r, []string{}, func(session interfaces.Session) {
		var out interface{}
		var err error
		switch strings.TrimPrefix(functionName, "sdk.activitypub.") {
		case "activity.create":
			var req rpc_net_hereus_sdk_activitypub_activity.CreateArguments
			if tryReadBody(r, w, &req) {
				out, err = rpc_net_hereus_sdk_activitypub_activity.Create(session, &req)
			}
		case "activity.announce":
			var req rpc_net_hereus_sdk_activitypub_activity.AnnounceArguments
			if tryReadBody(r, w, &req) {
				err = rpc_net_hereus_sdk_activitypub_activity.Announce(session, &req)
			}
		case "activity.undoAnnounce":
			var req rpc_net_hereus_sdk_activitypub_activity.UndoAnnounceArguments
			if tryReadBody(r, w, &req) {
				err = rpc_net_hereus_sdk_activitypub_activity.UndoAnnounce(session, &req)
			}
		case "activity.delete":
			var req rpc_net_hereus_sdk_activitypub_activity.DeleteArguments
			if tryReadBody(r, w, &req) {
				err = rpc_net_hereus_sdk_activitypub_activity.Delete(session, &req)
			}
		case "activity.get":
			var req rpc_net_hereus_sdk_activitypub_activity.GetArguments
			if tryReadBody(r, w, &req) {
				out, err = rpc_net_hereus_sdk_activitypub_activity.Get(session, &req)
			}
		case "activity.list":
			var req rpc_net_hereus_sdk_activitypub_activity.ListArguments
			if tryReadBody(r, w, &req) {
				out, err = rpc_net_hereus_sdk_activitypub_activity.List(session, &req)
			}
		case "activity.like":
			var req rpc_net_hereus_sdk_activitypub_activity.LikeArguments
			if tryReadBody(r, w, &req) {
				err = rpc_net_hereus_sdk_activitypub_activity.Like(session, &req)
			}
		case "activity.unlike":
			sdkaphandlers.HandleUnlike(kernel, w, r)
		case "activity.edit":
			sdkaphandlers.HandleEdit(kernel, w, r)
		case "actor.get":
			sdkaphandlers.HandleActor(kernel, w, r)
		case "actor.follow":
			var req rpc_net_hereus_sdk_activitypub_actor.FollowRequest
			if tryReadBody(r, w, &req) {
				err = rpc_net_hereus_sdk_activitypub_actor.Follow(session, &req)
			}
		case "actor.unfollow":
			var req rpc_net_hereus_sdk_activitypub_actor.UnfollowArguments
			if tryReadBody(r, w, &req) {
				err = rpc_net_hereus_sdk_activitypub_actor.Unfollow(session, &req)
			}
		case "actor.followers":
			sdkaphandlers.HandleFollowers(kernel, w, r)
		case "actor.following":
			sdkaphandlers.HandleFollowing(kernel, w, r)
		default:
			http.Error(w, fmt.Sprintf("Unknown endpoint: %s", endpoint), http.StatusNotFound)
		}
		convertRpcResponseToHttpResponse(out, err, w)
	})
}
