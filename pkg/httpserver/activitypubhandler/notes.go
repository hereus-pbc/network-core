package activitypubhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func HandleNotes(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request, activityId string) {
	if strings.Split(r.Header["Accept"][0], ",")[0] != "application/activity+json" {
		http.Redirect(w, r, fmt.Sprintf("https://social.hereus.net/statuses/hereus:%s", activityId), http.StatusSeeOther)
		return
	}
	note, err := kernel.ActivityPubDB().FetchNoteById("", activityId)
	if err != nil {
		http.Error(w, "Note not found or not public", http.StatusNotFound)
		return
	}
	if note.LdContext == nil {
		note.LdContext = "https://www.w3.org/ns/activitystreams"
	}
	w.Header().Set("Content-Type", "application/activity+json")
	if json.NewEncoder(w).Encode(note) != nil {
		http.Error(w, "Failed to serialize activity", http.StatusInternalServerError)
		return
	}
}
