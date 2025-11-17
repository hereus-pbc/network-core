package rpc_net_hereus_sdk_activitypub_activity

import (
	"fmt"

	"github.com/hereus-pbc/golang-utils/randomizer"
	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

func UndoAnnounceRpc() *helpers.RpcFunctionHandlerWithArguments {
	return &helpers.RpcFunctionHandlerWithArguments{
		ReqFactory: func() interface{} { return &UndoAnnounceArguments{} },
		Handler: func(session interfaces.Session, req interface{}) (interface{}, error) {
			return nil, UndoAnnounce(session, req.(*UndoAnnounceArguments))
		},
		Permissions: []string{"net.hereus.sdk.permissions.activitypub"},
	}
}

type UndoAnnounceArguments struct {
	ObjectId string `json:"objectId"` // ID of the object being announced ("object")
}

func UndoAnnounce(session interfaces.Session, req *UndoAnnounceArguments) error {
	announcements, err := session.GetKernel().ActivityPubDB().RemoveAnnounce(session.GetUser().GetActorUrl(), req.ObjectId)
	if err != nil {
		if err.Error() == "no announcements found" {
			return fmt.Errorf("no announcements found")
		}
		return fmt.Errorf("failed to undo announce: %w", err)
	}
	for _, announcement := range *announcements {
		uuid := randomizer.Random128ByteString()
		session.GetKernel().PushOutgoingActivity(types.ActivityStream{
			LdContext: "https://www.w3.org/ns/activitystreams",
			Id:        fmt.Sprintf("https://%s/activitypub/activities/undo-announce-%s", session.GetKernel().GetDomain(), uuid),
			Type:      "Undo",
			Actor:     session.GetUser().GetActorUrl(),
			Object: types.ActivityStream{
				LdContext: "https://www.w3.org/ns/activitystreams",
				Id:        announcement.Id,
				Type:      "Announce",
				Actor:     announcement.Who,
				Object:    announcement.What,
				To:        announcement.To,
				Cc:        announcement.Cc,
			},
			To: announcement.To,
			Cc: announcement.Cc,
		})
	}
	return nil
}
