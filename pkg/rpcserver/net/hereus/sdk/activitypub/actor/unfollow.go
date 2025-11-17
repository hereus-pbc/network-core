package rpc_net_hereus_sdk_activitypub_actor

import (
	"fmt"
	"strings"

	"github.com/hereus-pbc/golang-utils/randomizer"
	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

func UnfollowRpc() *helpers.RpcFunctionHandlerWithArguments {
	return &helpers.RpcFunctionHandlerWithArguments{
		ReqFactory: func() interface{} {
			return &UnfollowArguments{}
		},
		Handler: func(session interfaces.Session, req interface{}) (interface{}, error) {
			return nil, Unfollow(session, req.(*UnfollowArguments))
		},
		Permissions: []string{"net.hereus.sdk.permissions.activitypub"},
	}
}

type UnfollowArguments struct {
	Handle string `json:"handle"`
}

func Unfollow(session interfaces.Session, req *UnfollowArguments) error {
	if req.Handle == "" {
		return fmt.Errorf("handle is required")
	}
	if !strings.HasPrefix(req.Handle, "https://") {
		_, toI, err := session.GetKernel().GetActorUrlByHandler(req.Handle)
		if err != nil {
			return err
		}
		req.Handle = toI
	}
	results, err := session.GetKernel().ActivityPubDB().RemoveFollow(session.GetUser().GetActorUrl(), req.Handle)
	if err != nil {
		if err.Error() == "follow does not exist" {
			return nil
		}
		return err
	}
	for _, result := range *results {
		session.GetKernel().PushOutgoingActivity(types.ActivityStream{
			LdContext: "https://www.w3.org/ns/activitystreams",
			Id:        fmt.Sprintf("https://%s/activitypub/activities/unfollow-%s", session.GetKernel().GetDomain(), randomizer.Random128ByteString()),
			Type:      "Undo",
			Actor:     session.GetUser().GetActorUrl(),
			Object: types.ActivityStream{
				Id:     result.Id,
				Type:   "Follow",
				Actor:  session.GetUser().GetActorUrl(),
				Object: result.Whom,
				To:     result.Whom,
				Cc:     result.Whom,
			},
			To: result.Whom,
			Cc: result.Whom,
		})
	}
	return nil
}
