package rpc_net_hereus_sdk_activitypub_activity

import (
	"fmt"

	"github.com/hereus-pbc/golang-utils/randomizer"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

type UnlikeArguments struct {
	ObjectId string `json:"objectId"`
}

func Unlike(session interfaces.Session, req *UnlikeArguments) error {
	if req.ObjectId == "" {
		return fmt.Errorf("objectId is required")
	}
	note, err := session.GetKernel().ActivityPubDB().FetchNoteByUrl(session.GetUser().GetActorUrl(), req.ObjectId)
	if err != nil {
		return fmt.Errorf("failed to fetch note: %w", err)
	}
	id := fmt.Sprintf("https://%s/activitypub/activities/like-%s", session.GetKernel().GetDomain(), randomizer.Random128ByteString())
	ids, err := session.GetKernel().ActivityPubDB().UnlikeNote(session.GetUser().GetActorUrl(), note.Id)
	if err != nil {
		return fmt.Errorf("failed to unlike note: %w", err)
	}
	for _, likeId := range ids {
		session.GetKernel().PushOutgoingActivity(types.ActivityStream{
			LdContext: "https://www.w3.org/ns/activitystreams",
			Id:        id,
			Type:      "Undo",
			Actor:     session.GetUser().GetActorUrl(),
			To:        "https://www.w3.org/ns/activitystreams#Public",
			Cc:        []string{note.AttributedTo},
			Object: types.ActivityStream{
				LdContext: "https://www.w3.org/ns/activitystreams",
				Id:        likeId,
				Type:      "Like",
				Actor:     session.GetUser().GetActorUrl(),
				Object:    note.Id,
				To:        "https://www.w3.org/ns/activitystreams#Public",
				Cc:        []string{note.AttributedTo},
			},
		})
	}
	return nil
}
