package activitypubhandler

import (
	"encoding/json"
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func HandleActivities(kernel interfaces.Kernel, w http.ResponseWriter, _ *http.Request, activityId string) {
	activity, err := kernel.ActivityPubDB().FetchActivityById("", activityId)
	if err != nil {
		http.Error(w, "Activity not found or not public", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/activity+json")
	if json.NewEncoder(w).Encode(activity) != nil {
		http.Error(w, "Failed to serialize activity", http.StatusInternalServerError)
		return
	}
}
