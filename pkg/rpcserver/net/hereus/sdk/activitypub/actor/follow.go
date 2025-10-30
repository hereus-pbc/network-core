package rpc_net_hereus_sdk_activitypub_actor

import (
	"fmt"
	"strings"
	"time"

	"github.com/hereus-pbc/golang-utils/randomizer"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

type FollowRequest struct {
	Handle string `json:"handle"`
}

func Follow(session interfaces.Session, req *FollowRequest) error {
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
	id := fmt.Sprintf("https://%s/activitypub/activities/follow-%s", session.GetKernel().GetDomain(), randomizer.Random128ByteString())
	now := time.Now().UTC()
	err := session.GetKernel().ActivityPubDB().AddFollow(session.GetUser().GetActorUrl(), req.Handle, now.Format(time.DateTime), id)
	if err != nil {
		return err
	}
	session.GetKernel().PushOutgoingActivity(types.ActivityStream{
		LdContext: "https://www.w3.org/ns/activitystreams",
		Id:        id,
		Type:      "Follow",
		Actor:     session.GetUser().GetActorUrl(),
		Object:    req.Handle,
		To:        req.Handle,
	})
	return nil
}
