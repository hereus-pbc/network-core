package activitypubhandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hereus-pbc/golang-utils/cryptography"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

func HandleInboxes(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request) {
	var activityMap map[string]interface{}
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body")
	}
	err = json.Unmarshal(rBody, &activityMap)
	if err != nil {
		fmt.Println("request body is not a valid JSON")
	}
	actor, err := kernel.ActivityPubDB().GetActor(activityMap["actor"].(string), "")
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Error fetching actor: %s\n", err)
		http.Error(w, "Error fetching actor", http.StatusBadRequest)
		return
	}
	err = cryptography.VerifyHttpSignature(r, r.Header.Get("Signature"), actor.PublicKey.PublicKey, kernel.GetDomain())
	if err != nil {
		fmt.Println("Signature verification failed")
		fmt.Println(err)
		http.Error(w, "Signature verification failed", http.StatusUnauthorized)
		return
	}
	activity := types.ActivityStream{
		LdContext: "https://www.w3.org/ns/activitystreams",
		Id:        activityMap["id"].(string),
		Type:      activityMap["type"].(string),
		Actor:     activityMap["actor"].(string),
		Object:    nil,
		To:        nil,
		Cc:        make([]string, 0),
	}

	if to, ok := activityMap["to"].(string); ok {
		activity.To = to
	} else if toArray, ok := activityMap["to"].([]interface{}); ok {
		activity.To = toArray[0].(string)
		for _, to := range toArray {
			if toStr, ok := to.(string); ok {
				activity.Cc = append(activity.Cc.([]string), toStr)
			}
		}
	}
	if cc, ok := activityMap["cc"].(string); ok {
		activity.Cc = append(activity.Cc.([]string), cc)
	} else if ccArray, ok := activityMap["cc"].([]interface{}); ok {
		for _, cc := range ccArray {
			if ccStr, ok := cc.(string); ok {
				activity.Cc = append(activity.Cc.([]string), ccStr)
			}
		}
	}
	switch activityMap["object"].(type) {
	case map[string]interface{}:
		switch activityMap["object"].(map[string]interface{})["type"].(string) {
		case "Note":
			var note types.ActivityPubNote
			noteBytes, err := json.Marshal(activityMap["object"])
			if err != nil {
				fmt.Println("Error marshalling activity object")
				http.Error(w, "Activity object is not a valid type (Note1)", http.StatusBadRequest)
				return
			}
			if json.Unmarshal(noteBytes, &note) != nil {
				fmt.Println("Error unmarshalling activity object")
				http.Error(w, "Activity object is not a valid type (Note2)", http.StatusBadRequest)
				return
			}
			activity.Object = note
		case "Person":
			var actor types.Actor
			actorBytes, err := json.Marshal(activityMap["object"])
			if err != nil {
				fmt.Println("Error marshalling activity object")
				http.Error(w, "Activity object is not a valid type (Person1)", http.StatusBadRequest)
				return
			}
			if json.Unmarshal(actorBytes, &actor) != nil {
				fmt.Println("Error unmarshalling activity object")
				http.Error(w, "Activity object is not a valid type (Person2)", http.StatusBadRequest)
				return
			}
			activity.Object = actor
		default:
			var activityObject types.ActivityStream
			objectBytes, err := json.Marshal(activityMap["object"])
			if err != nil {
				fmt.Println("Error marshalling activity object")
				http.Error(w, "Activity object is not a valid type (Activity1)", http.StatusBadRequest)
				return
			}
			if json.Unmarshal(objectBytes, &activityObject) != nil {
				fmt.Println("Error unmarshalling activity object")
				http.Error(w, "Activity object is not a valid type (Activity2)", http.StatusBadRequest)
				return
			}
			activity.Object = activityObject
		}
	case string:
		activity.Object = activityMap["object"].(string)
	default:
		activity.Object = nil
	}

	// check if object is string
	if activity.Object == nil {
		fmt.Println("Activity object is nil")
		http.Error(w, "Activity object is nil", http.StatusBadRequest)
		return
	}
	kernel.PushIncomingActivity(activity)
	w.WriteHeader(200)
}
