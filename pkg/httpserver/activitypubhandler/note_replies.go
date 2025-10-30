package activitypubhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

type RepliesCollectionPage struct {
	LdContext string                  `json:"@context,omitempty"`
	Id        string                  `json:"id"`
	Type      string                  `json:"type"`
	PartOf    string                  `json:"partOf"`
	Items     []types.ActivityPubNote `json:"items"`
}

type RepliesCollection struct {
	LdContext  string                `json:"@context"`
	Id         string                `json:"id"`
	Type       string                `json:"type"`
	First      RepliesCollectionPage `json:"first"`
	TotalItems int                   `json:"totalItems"`
}

func HandleNoteReplies(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request, activityId string) {
	_, err := kernel.ActivityPubDB().FetchNoteById("", activityId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	replies := kernel.ActivityPubDB().ListRepliesToNote("", activityId)
	for i := range replies {
		if replies[i].LdContext == nil {
			replies[i].LdContext = "https://www.w3.org/ns/activitystreams"
		}
		replies[i].CustomUrl = ""
	}
	collectionPage := RepliesCollectionPage{
		LdContext: "https://www.w3.org/ns/activitystreams",
		Id:        fmt.Sprintf("https://%s/activitypub/note-replies/%s?page=1", kernel.GetDomain(), activityId),
		Type:      "CollectionPage",
		PartOf:    fmt.Sprintf("https://%s/activitypub/note-replies/%s", kernel.GetDomain(), activityId),
		Items:     replies,
	}
	if r.URL.Query().Get("page") == "true" || r.URL.Query().Get("page") == "1" {
		w.Header().Set("Content-Type", "application/activity+json")
		if json.NewEncoder(w).Encode(collectionPage) != nil {
			http.Error(w, "Failed to serialize activity", http.StatusInternalServerError)
		}
		return
	}
	collectionResponse := RepliesCollection{
		LdContext:  "https://www.w3.org/ns/activitystreams",
		Id:         fmt.Sprintf("https://%s/activitypub/note-replies/%s", kernel.GetDomain(), activityId),
		Type:       "Collection",
		First:      collectionPage,
		TotalItems: len(replies),
	}
	w.Header().Set("Content-Type", "application/activity+json")
	if json.NewEncoder(w).Encode(collectionResponse) != nil {
		http.Error(w, "Failed to serialize activity", http.StatusInternalServerError)
		return
	}
}
