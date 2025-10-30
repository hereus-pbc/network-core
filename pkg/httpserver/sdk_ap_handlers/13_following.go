package sdkaphandlers

import (
	"encoding/json"
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func HandleFollowing(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request) {
	kernel.SessionWrapper(w, r, []string{}, func(session interfaces.Session) {
		var resp []HandleActorResponse
		for _, followingUrl := range session.GetUser().ListFollowing() {
			following, err := kernel.ActivityPubDB().GetActor(followingUrl, "")
			if err != nil {
				continue
			}
			followingHandle, err := kernel.ReverseActorUrl(followingUrl, following)
			if err == nil {
				resp = append(resp, ActorToResponse(following, followingHandle))
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
