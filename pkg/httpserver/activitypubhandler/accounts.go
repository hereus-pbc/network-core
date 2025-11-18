package activitypubhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func HandleAccounts(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request, username string) {
	acceptMimetype := r.Header.Get("Accept")
	acceptMimetype = strings.SplitN(acceptMimetype, ";", 2)[0]
	acceptMimetypes := strings.Split(acceptMimetype, ",")
	if len(acceptMimetypes) != 0 &&
		acceptMimetypes[0] != "application/activity+json" && acceptMimetypes[0] != "application/ld+json" {
		if len(acceptMimetypes) == 2 &&
			acceptMimetypes[1] != "application/activity+json" && acceptMimetypes[1] != "application/ld+json" {
			http.Redirect(w, r, fmt.Sprintf("/~%s", username), http.StatusSeeOther)
			return
		}
	}
	user, err := kernel.UserManager().Get(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	actor, err := user.BuildActorObject()
	if err != nil {
		http.Error(w, "Failed to build actor object", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/activity+json")
	if json.NewEncoder(w).Encode(actor) != nil {
		http.Error(w, "Failed to encode actor", http.StatusInternalServerError)
		return
	}
}
