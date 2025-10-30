package sdkaphandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hereus-pbc/golang-utils/randomizer"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

type HandleEditRequest struct {
	ObjectId       string      `json:"objectId"`       // ID of the object being edited ("id")
	Summary        interface{} `json:"summary"`        // Plaintext summary of the activity ("content")
	ContentWarning interface{} `json:"contentWarning"` // Content warning ("summary", "sensitive" => if != "")
	Attachments    interface{} `json:"attachments"`    // Attachment URLs ("attachment", optional)
	Properties     interface{} `json:"properties"`     // Additional custom properties
	Url            interface{} `json:"url"`            // URL of the activity (optional, server will generate if empty)
}

func SelfIfNotNilElseOther(a interface{}, b interface{}) interface{} {
	if a != nil {
		return a
	}
	return b
}

func HandleEdit(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request) {
	kernel.SessionWrapper(w, r, []string{}, func(session interfaces.Session) {
		var req HandleEditRequest
		if json.NewDecoder(r.Body).Decode(&req) != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		note, err := kernel.ActivityPubDB().FetchNoteByUrl(session.GetUser().GetActorUrl(), req.ObjectId)
		if err != nil {
			http.Error(w, "Failed to fetch object: "+err.Error(), http.StatusBadRequest)
			return
		}
		if note.AttributedTo != session.GetUser().GetActorUrl() {
			http.Error(w, "You can only edit your own objects", http.StatusForbidden)
			return
		}
		reqSensitive := note.Sensitive
		if req.ContentWarning != nil {
			if str, ok := req.ContentWarning.(string); ok {
				reqSensitive = str != ""
			}
		}
		err = kernel.ActivityPubDB().EditNote(types.ActivityPubNote{
			LdContext:    note.LdContext,
			Id:           note.Id,
			Type:         "Note",
			Summary:      SelfIfNotNilElseOther(req.ContentWarning, note.Summary).(string),
			InReplyTo:    note.InReplyTo,
			Published:    note.Published,
			Updated:      time.Now().UTC().Format(time.RFC3339),
			Url:          note.Url,
			CustomUrl:    SelfIfNotNilElseOther(req.Url, note.CustomUrl).(string),
			AttributedTo: note.AttributedTo,
			To:           note.To,
			Cc:           note.Cc,
			Sensitive:    reqSensitive,
			Content:      SelfIfNotNilElseOther(req.Summary, note.Content).(string),
			Tag:          note.Tag,
			Replies:      note.Replies,
			Properties:   SelfIfNotNilElseOther(req.Properties, note.Properties).(map[string]interface{}),
			PackageName:  "",
			Attachment: (func() []types.ActivityPubAttachment {
				if req.Attachments == nil {
					return note.Attachment
				}
				var attachments []types.ActivityPubAttachment
				arr, ok := req.Attachments.([]string)
				if !ok {
					return note.Attachment
				}
				for _, url := range arr {
					attachments = append(attachments, types.ActivityPubAttachment{
						Type:      "Document",
						MediaType: "application/octet-stream",
						Url:       url,
					})
				}
				return attachments
			})(),
		})
		if err != nil {
			http.Error(w, "Failed to update object: "+err.Error(), http.StatusInternalServerError)
			return
		}
		newNote, err := kernel.ActivityPubDB().FetchNoteByUrl(session.GetUser().GetActorUrl(), note.Id)
		if err != nil {
			http.Error(w, "Failed to fetch updated object: "+err.Error(), http.StatusInternalServerError)
			return
		}
		kernel.PushOutgoingActivity(types.ActivityStream{
			LdContext: "https://www.w3.org/ns/activitystreams",
			Id:        fmt.Sprintf("https://%s/activitypub/activities/edit-%s", kernel.GetDomain(), randomizer.Random128ByteString()),
			Type:      "Update",
			Actor:     session.GetUser().GetActorUrl(),
			Object:    newNote,
			To:        newNote.To,
			Cc: (func() []string {
				if ccList, ok := newNote.Cc.([]string); ok {
					return ccList
				}
				return []string{}
			})(),
		})
	})
}
