package sdkaphandlers

import (
	"encoding/json"
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func HandleFollowers(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request) {
	kernel.SessionWrapper(w, r, []string{}, func(session interfaces.Session) {
		var resp []HandleActorResponse
		for _, followerUrl := range session.GetUser().ListFollowers() {
			follower, err := kernel.ActivityPubDB().GetActor(followerUrl, "")
			if err != nil {
				continue
			}
			followerHandle, err := kernel.ReverseActorUrl(followerUrl, follower)
			if err == nil {
				resp = append(resp, ActorToResponse(follower, followerHandle))
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if len(resp) == 0 {
			w.Write([]byte("[]"))
			return
		}
		json.NewEncoder(w).Encode(resp)
	})
}
