package rpc_net_hereus_sdk_activitypub_activity

import (
	"fmt"

	"github.com/hereus-pbc/golang-utils/randomizer"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

type AnnounceArguments struct {
	ObjectId string   `json:"objectId"` // ID of the object being announced ("object")
	To       string   `json:"to"`       // Audience ("to")
	Cc       []string `json:"cc"`       // Audience ("cc", optional)
}

func Announce(session interfaces.Session, req *AnnounceArguments) error {
	if req.ObjectId == "" || req.To == "" {
		return fmt.Errorf("missing required fields")
	}
	if session.GetKernel().ActivityPubDB().HasAnnounced(session.GetUser().GetActorUrl(), req.ObjectId) {
		return fmt.Errorf("object already announced")
	}
	note, err := session.GetKernel().ActivityPubDB().FetchNoteByUrl(session.GetUser().GetActorUrl(), req.ObjectId)
	if err != nil {
		return fmt.Errorf("failed to fetch object: %w", err)
	}
	if req.Cc == nil {
		req.Cc = []string{}
	}
	CcMap := map[string]bool{}
	for _, cc := range req.Cc {
		CcMap[cc] = true
	}
	if !CcMap[note.AttributedTo] {
		req.Cc = append(req.Cc, note.AttributedTo)
	}
	uuid := randomizer.Random128ByteString()
	session.GetKernel().PushOutgoingActivity(types.ActivityStream{
		LdContext: "https://www.w3.org/ns/activitystreams",
		Id:        fmt.Sprintf("https://%s/activitypub/activities/announce-%s", session.GetKernel().GetDomain(), uuid),
		Type:      "Announce",
		Actor:     session.GetUser().GetActorUrl(),
		Object:    req.ObjectId,
		To:        req.To,
		Cc:        req.Cc,
	})
	return nil
}
