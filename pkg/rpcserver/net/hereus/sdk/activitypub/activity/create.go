package rpc_net_hereus_sdk_activitypub_activity

import (
	"fmt"
	"time"

	"github.com/hereus-pbc/golang-utils/randomizer"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/misc/texttohtml"
	"github.com/hereus-pbc/network-core/pkg/types"
)

type ActivityType struct {
	Summary        string                 `json:"summary"`        // Plaintext summary of the activity ("content")
	ContentWarning string                 `json:"contentWarning"` // Content warning ("summary", "sensitive" => if != "")
	InReplyTo      string                 `json:"inReplyTo"`      // ID of the object being replied to, if any ("inReplyTo")
	To             string                 `json:"to"`             // Audience ("to")
	Cc             []string               `json:"cc"`             // Audience ("cc", optional)
	Attachments    []string               `json:"attachments"`    // Attachment URLs ("attachment", optional)
	Properties     map[string]interface{} `json:"properties"`     // Additional custom properties
	Url            string                 `json:"url"`            // URL of the activity (optional, server will generate if empty)
	Owner          string                 `json:"owner"`          // Owner of the activity
	Published      string                 `json:"published"`      // Publication date (optional, server will generate if empty)
	Id             string                 `json:"id"`             // ID of the activity (optional, server will generate if empty)
}

type CreateArguments struct {
	Summary        string                 `json:"summary"`        // Plaintext summary of the activity ("content")
	ContentWarning string                 `json:"contentWarning"` // Content warning ("summary", "sensitive" => if != "")
	InReplyTo      string                 `json:"inReplyTo"`      // ID of the object being replied to, if any ("inReplyTo")
	To             string                 `json:"to"`             // Audience ("to")
	Cc             []string               `json:"cc"`             // Audience ("cc", optional)
	Attachments    []string               `json:"attachments"`    // Attachment URLs ("attachment", optional)
	Properties     map[string]interface{} `json:"properties"`     // Additional custom properties
	Url            string                 `json:"url"`            // URL of the activity (optional, server will generate if empty)
}

func Create(session interfaces.Session, req *CreateArguments) (string, error) {
	if req.To == "" {
		return "", fmt.Errorf("invalid 'to' field")
	}
	toL, err := texttohtml.ConvertAllHandlesToUrls(session.GetKernel(), req.To)
	if err != nil {
		return "", fmt.Errorf("invalid 'to' field: %w", err)
	}
	to := toL[req.To]

	if req.Cc == nil {
		req.Cc = []string{}
	}
	ccMap, err := texttohtml.ConvertAllHandlesToUrls(session.GetKernel(), req.Cc)
	if err != nil {
		return "", fmt.Errorf("invalid 'cc' field: %w", err)
	}
	var cc []string
	for _, ccUrl := range ccMap {
		cc = append(cc, ccUrl)
	}

	content, hashtags, mentions := texttohtml.TextToHtml(session.GetKernel(), session.GetKernel().GetDomain(), req.Summary)

	uuid := randomizer.Random128ByteString()
	noteTime := time.Now().UTC()
	urlId := fmt.Sprintf("https://%s/activitypub/notes/%s", session.GetKernel().GetDomain(), uuid)
	session.GetKernel().PushOutgoingActivity(types.ActivityStream{
		LdContext: "https://www.w3.org/ns/activitystreams",
		Id:        fmt.Sprintf("https://%s/activitypub/activities/create-%s", session.GetKernel().GetDomain(), uuid),
		Type:      "Create",
		Actor:     session.GetUser().GetActorUrl(),
		To:        to,
		Cc:        cc,
		Object: types.ActivityPubNote{
			Id:           urlId,
			Type:         "Note",
			InReplyTo:    req.InReplyTo,
			Published:    noteTime.Format(time.RFC3339),
			Url:          fmt.Sprintf("https://%s/activitypub/notes/%s", session.GetKernel().GetDomain(), uuid),
			AttributedTo: session.GetUser().GetActorUrl(),
			To:           to,
			Cc:           cc,
			Sensitive:    req.ContentWarning != "",
			Content:      content,
			Replies:      fmt.Sprintf("https://%s/activitypub/note-replies/%s", session.GetKernel().GetDomain(), uuid),
			PackageName:  session.GetSession().PackageName,
			Tag: (func() []types.ActivityPubTag {
				tags := []types.ActivityPubTag{
					{
						Type: "Mention",
						Name: req.To,
						Href: to,
					},
				}
				for _, c := range req.Cc {
					tags = append(tags, types.ActivityPubTag{
						Type: "Mention",
						Name: c,
						Href: ccMap[c],
					})
				}
				for _, hashtag := range hashtags {
					tags = append(tags, types.ActivityPubTag{
						Type: "Hashtag",
						Name: "#" + hashtag,
						Href: fmt.Sprintf("https://%s/tags/%s", session.GetKernel().GetDomain(), hashtag),
					})
				}
				for handle, url := range mentions {
					tags = append(tags, types.ActivityPubTag{
						Type: "Mention",
						Name: handle,
						Href: url,
					})
				}
				return tags
			})(),
			Attachment: (func() []types.ActivityPubAttachment {
				var attachments []types.ActivityPubAttachment
				for _, url := range req.Attachments {
					attachments = append(attachments, types.ActivityPubAttachment{
						Type:      "Document",
						MediaType: "application/octet-stream",
						Url:       url,
					})
				}
				return attachments
			})(),
			Summary: (func() string {
				if req.ContentWarning != "" {
					return req.ContentWarning
				} else {
					return ""
				}
			})(),
			Properties: req.Properties,
		},
	})
	return urlId, nil
}
