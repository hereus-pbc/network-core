package rpc_net_hereus_sdk_activitypub_activity

import (
	"strings"
	"time"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

func NoteToSdkObject(kernel interfaces.Kernel, note types.ActivityPubNote) (*ActivityType, error) {
	actor, err := kernel.ActivityPubDB().GetActor(note.AttributedTo, "")
	if err != nil {
		return nil, err
	}
	owner, err := kernel.ReverseActorUrl(note.AttributedTo, actor)
	if err != nil {
		return nil, err
	}
	owner = strings.TrimPrefix(owner, "@")
	var to string
	var cc []string
	if note.To == "https://www.w3.org/ns/activitystreams#Public" {
		to = "Public"
	} else if arr, ok := note.To.([]string); ok {
		to = arr[0]
		for _, v := range arr {
			if v == "https://www.w3.org/ns/activitystreams#Public" {
				to = "Public"
			} else {
				cc = append(cc, v)
			}
		}
	}
	if arr, ok := note.Cc.([]string); ok {
		for _, v := range arr {
			if v == "https://www.w3.org/ns/activitystreams#Public" {
				to = "Public"
			} else if v != to {
				cc = append(cc, v)
			}
		}
	}
	published, err := time.Parse(time.RFC3339, note.Published)
	if err != nil {
		published = time.Now()
	}
	return &ActivityType{
		Summary:        note.Content,
		ContentWarning: note.Summary,
		InReplyTo:      note.InReplyTo,
		Url:            note.Url,
		Owner:          owner,
		To:             to,
		Cc:             cc,
		Published:      published.Format(time.DateTime),
		Id:             note.Id,
		Properties: (func() map[string]interface{} {
			if note.Properties == nil {
				return map[string]interface{}{}
			} else {
				return note.Properties
			}
		})(),
		Attachments: (func() []string {
			return []string{}
		})(),
	}, nil
}

type GetArguments struct {
	ObjectId string `json:"objectId"` // ID of the object being requested ("object")
}

func Get(session interfaces.Session, req *GetArguments) (*ActivityType, error) {
	note, err := session.GetKernel().ActivityPubDB().FetchNoteByUrl(session.GetUser().GetActorUrl(), req.ObjectId)
	if err != nil {
		return nil, err
	}
	return NoteToSdkObject(session.GetKernel(), note)
}
