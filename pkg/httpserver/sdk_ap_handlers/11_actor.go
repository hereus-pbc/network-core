package sdkaphandlers

import (
	"encoding/json"
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

type HandleActorRequest struct {
	Handle string `json:"handle"` // URL of the actor to handle
}

type HandleActorResponse struct {
	Handle         string `json:"handle"`
	Url            string `json:"url"`
	DisplayName    string `json:"displayName"`
	Biography      string `json:"biography"`
	ProfilePicture string `json:"profilePicture"`
	BannerPicture  string `json:"bannerPicture"`
}

func ActorToResponse(actor types.Actor, handle string) HandleActorResponse {
	return HandleActorResponse{
		Handle:         handle,
		Url:            actor.Id,
		DisplayName:    actor.Name,
		Biography:      actor.Summary,
		ProfilePicture: actor.Icon.Url,
		BannerPicture:  actor.Image.Url,
	}
}

func HandleActor(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request) {
	kernel.SessionWrapper(w, r, []string{}, func(session interfaces.Session) {
		var req HandleActorRequest
		if json.NewDecoder(r.Body).Decode(&req) != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		_, url, err := kernel.GetActorUrlByHandler(req.Handle)
		if err != nil {
			http.Error(w, "Failed to fetch actor: "+err.Error(), http.StatusBadRequest)
			return
		}
		actor, err := kernel.ActivityPubDB().GetActor(url, "")
		if err != nil {
			http.Error(w, "Failed to fetch actor: "+err.Error(), http.StatusBadRequest)
			return
		}
		resp := ActorToResponse(actor, req.Handle)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}
