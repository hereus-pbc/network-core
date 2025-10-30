package http_server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hereus-pbc/network-core/pkg/httpserver/activitypubhandler"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func dispatchActivityPub(kernel interfaces.Kernel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint, _ := strings.CutPrefix(r.URL.Path, "/activitypub/")
		endpoint = strings.TrimSuffix(endpoint, "/")
		endpoint = strings.ReplaceAll(endpoint, "//", "/")
		if len(endpoint) == 0 {
			http.Error(w, "Endpoint not specified", http.StatusBadRequest)
			return
		}
		fmt.Printf("(ActivityPub) %s [%s] %s\n", r.RemoteAddr, time.Now().Format("2006-01-02 15:04:05"), endpoint)
		endpointPieces := strings.Split(endpoint, "/")
		if len(endpointPieces) < 1 {
			http.Error(w, "Invalid endpoint", http.StatusBadRequest)
			return
		}
		switch endpointPieces[0] {
		case "accounts":
			activitypubhandler.HandleAccounts(kernel, w, r, endpointPieces[1])
		case "inboxes":
			activitypubhandler.HandleInboxes(kernel, w, r)
		case "outboxes":
			activitypubhandler.HandleCollection(kernel, w, r, "outboxes", endpointPieces[1])
		case "followers":
			activitypubhandler.HandleCollection(kernel, w, r, "followers", endpointPieces[1])
		case "following":
			activitypubhandler.HandleCollection(kernel, w, r, "following", endpointPieces[1])
		case "inbox":
			activitypubhandler.HandleInboxes(kernel, w, r)
		case "activities":
			activitypubhandler.HandleActivities(kernel, w, r, endpointPieces[1])
		case "notes":
			activitypubhandler.HandleNotes(kernel, w, r, endpointPieces[1])
		case "note-replies":
			activitypubhandler.HandleNoteReplies(kernel, w, r, endpointPieces[1])
		default:
			http.Error(w, fmt.Sprintf("Reference not found: %s", endpoint), http.StatusNotFound)
		}
	}
}
