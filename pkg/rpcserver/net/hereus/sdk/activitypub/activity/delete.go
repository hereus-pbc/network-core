package rpc_net_hereus_sdk_activitypub_activity

import (
	"fmt"

	"github.com/hereus-pbc/golang-utils/randomizer"
	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

func DeleteRpc() *helpers.RpcFunctionHandlerWithArguments {
	return &helpers.RpcFunctionHandlerWithArguments{
		ReqFactory: func() interface{} { return &DeleteArguments{} },
		Handler: func(session interfaces.Session, req interface{}) (interface{}, error) {
			return nil, Delete(session, req.(*DeleteArguments))
		},
		Permissions: []string{"net.hereus.sdk.permissions.activitypub"},
	}
}

type DeleteArguments struct {
	ObjectId string `json:"objectId"`
}

func Delete(session interfaces.Session, req *DeleteArguments) error {
	if req.ObjectId == "" {
		return fmt.Errorf("objectId is required")
	}
	note, err := session.GetKernel().ActivityPubDB().FetchNoteById(session.GetUser().GetActorUrl(), req.ObjectId)
	if err != nil {
		return fmt.Errorf("failed to fetch note: %w", err)
	}
	if note.AttributedTo != session.GetUser().GetActorUrl() {
		return fmt.Errorf("you are not the author of this note")
	}
	var cc []string
	if note.Cc != nil {
		for _, v := range note.Cc.([]string) {
			cc = append(cc, v)
		}
	}
	cc = append(cc, note.To.(string))
	uuid := randomizer.Random128ByteString()
	session.GetKernel().PushOutgoingActivity(types.ActivityStream{
		LdContext: "https://www.w3.org/ns/activitystreams",
		Id:        fmt.Sprintf("https://%s/activitypub/activities/delete-%s", session.GetKernel().GetDomain(), uuid),
		Type:      "Delete",
		Actor:     session.GetUser().GetActorUrl(),
		Object:    note,
		To:        "https://www.w3.org/ns/activitystreams#Public",
		Cc:        cc,
	})
	return session.GetKernel().ActivityPubDB().DeleteNote(note.Id)
}
