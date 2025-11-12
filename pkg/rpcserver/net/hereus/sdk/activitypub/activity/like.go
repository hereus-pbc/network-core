package rpc_net_hereus_sdk_activitypub_activity

import (
	"fmt"
	"time"

	"github.com/hereus-pbc/golang-utils/randomizer"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

type LikeArguments struct {
	ObjectId string `json:"objectId"`
}

func Like(session interfaces.Session, req *LikeArguments) error {
	if req.ObjectId == "" {
		return fmt.Errorf("objectId is required")
	}
	note, err := session.GetKernel().ActivityPubDB().FetchNoteByUrl(session.GetUser().GetActorUrl(), req.ObjectId)
	if err != nil {
		return fmt.Errorf("failed to fetch note: %w", err)
	}
	id := fmt.Sprintf("https://%s/activitypub/activities/like-%s", session.GetKernel().GetDomain(), randomizer.Random128ByteString())
	now := time.Now()
	err = session.GetKernel().ActivityPubDB().RecordLike(session.GetUser().GetActorUrl(), note.Id, now.Format(time.DateTime), id)
	if err != nil {
		return fmt.Errorf("failed to like note: %w", err)
	}
	session.GetKernel().PushOutgoingActivity(types.ActivityStream{
		LdContext: "https://www.w3.org/ns/activitystreams",
		Id:        id,
		Type:      "Like",
		Actor:     session.GetUser().GetActorUrl(),
		Object:    note.Id,
		To:        "https://www.w3.org/ns/activitystreams#Public",
		Cc:        []string{note.AttributedTo},
	})
	return nil
}
