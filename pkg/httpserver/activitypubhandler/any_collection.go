package activitypubhandler

import (
	"encoding/json"
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

type CollectionResponse struct {
	LdContext  []string `json:"@context"`
	Id         string   `json:"id"`
	Type       string   `json:"type"`
	TotalItems int      `json:"totalItems"`
}

func HandleCollection(kernel interfaces.Kernel, w http.ResponseWriter, _ *http.Request, collectionName string, username string) {
	user, err := kernel.UserManager().Get(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	var totalItems int
	switch collectionName {
	case "outboxes":
		notes, err := kernel.ActivityPubDB().ListNotes("", map[string]interface{}{"owner": user.GetActorUrl()}, "", 0, false)
		if err != nil {
			http.Error(w, "Failed to list notes", http.StatusInternalServerError)
			return
		}
		totalItems = len(notes)
	case "followers":
		totalItems = len(user.ListFollowers())
	case "following":
		totalItems = len(user.ListFollowing())
	default:
		http.Error(w, "Unknown collection", http.StatusBadRequest)
		return
	}
	collectionId := user.GetRootActivityResourceUrl(collectionName)
	response := CollectionResponse{
		LdContext:  []string{"https://www.w3.org/ns/activitystreams"},
		Id:         collectionId,
		Type:       "OrderedCollection",
		TotalItems: totalItems,
	}
	w.Header().Set("Content-Type", "application/activity+json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
