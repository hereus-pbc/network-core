package sdkaphandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hereus-pbc/golang-utils/randomizer"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

type HandleUnlikeRequest struct {
	ObjectId string `json:"objectId"`
}

func HandleUnlike(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request) {
	kernel.SessionWrapper(w, r, []string{}, func(session interfaces.Session) {
		var req HandleUnlikeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.ObjectId == "" {
			http.Error(w, "Missing objectId", http.StatusBadRequest)
			return
		}
		note, err := kernel.ActivityPubDB().FetchNoteByUrl(session.GetUser().GetActorUrl(), req.ObjectId)
		if err != nil {
			http.Error(w, "Failed to fetch object: "+err.Error(), http.StatusBadRequest)
			return
		}
		id := fmt.Sprintf("https://%s/activitypub/activities/like-%s", kernel.GetDomain(), randomizer.Random128ByteString())
		ids, err := kernel.ActivityPubDB().UnlikeNote(session.GetUser().GetActorUrl(), note.Id)
		if err != nil {
			http.Error(w, "Failed to process undo: "+err.Error(), http.StatusInternalServerError)
			return
		}
		for _, likeId := range ids {
			kernel.PushOutgoingActivity(types.ActivityStream{
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
		w.WriteHeader(http.StatusOK)
	})
}
