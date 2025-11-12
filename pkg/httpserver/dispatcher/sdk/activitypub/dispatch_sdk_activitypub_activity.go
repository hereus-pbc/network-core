package dispatcher_sdk_activitypub

import (
	"fmt"
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	rpc_net_hereus_sdk_activitypub_activity "github.com/hereus-pbc/network-core/pkg/rpcserver/net/hereus/sdk/activitypub/activity"
)

func handleActivityPubSDKActivity(session interfaces.Session, w http.ResponseWriter, r *http.Request, innerFunction string) (out interface{}, err error) {
	switch innerFunction {
	case "create":
		var req rpc_net_hereus_sdk_activitypub_activity.CreateArguments
		if helpers.TryReadBody(r, w, &req) {
			return rpc_net_hereus_sdk_activitypub_activity.Create(session, &req)
		}
	case "announce":
		var req rpc_net_hereus_sdk_activitypub_activity.AnnounceArguments
		if helpers.TryReadBody(r, w, &req) {
			return nil, rpc_net_hereus_sdk_activitypub_activity.Announce(session, &req)
		}
	case "undoAnnounce":
		var req rpc_net_hereus_sdk_activitypub_activity.UndoAnnounceArguments
		if helpers.TryReadBody(r, w, &req) {
			return nil, rpc_net_hereus_sdk_activitypub_activity.UndoAnnounce(session, &req)
		}
	case "delete":
		var req rpc_net_hereus_sdk_activitypub_activity.DeleteArguments
		if helpers.TryReadBody(r, w, &req) {
			return nil, rpc_net_hereus_sdk_activitypub_activity.Delete(session, &req)
		}
	case "get":
		var req rpc_net_hereus_sdk_activitypub_activity.GetArguments
		if helpers.TryReadBody(r, w, &req) {
			return rpc_net_hereus_sdk_activitypub_activity.Get(session, &req)
		}
	case "list":
		var req rpc_net_hereus_sdk_activitypub_activity.ListArguments
		if helpers.TryReadBody(r, w, &req) {
			return rpc_net_hereus_sdk_activitypub_activity.List(session, &req)
		}
	case "like":
		var req rpc_net_hereus_sdk_activitypub_activity.LikeArguments
		if helpers.TryReadBody(r, w, &req) {
			return nil, rpc_net_hereus_sdk_activitypub_activity.Like(session, &req)
		}
	case "unlike":
		var req rpc_net_hereus_sdk_activitypub_activity.UnlikeArguments
		if helpers.TryReadBody(r, w, &req) {
			return nil, rpc_net_hereus_sdk_activitypub_activity.Unlike(session, &req)
		}
	case "edit":
		var req rpc_net_hereus_sdk_activitypub_activity.EditArguments
		if helpers.TryReadBody(r, w, &req) {
			return nil, rpc_net_hereus_sdk_activitypub_activity.Edit(session, &req)
		}
	default:
		return nil, fmt.Errorf("unknown activity endpoint: net.hereus.sdk.activitypub.activity.%s", innerFunction)
	}
	return nil, nil
}
