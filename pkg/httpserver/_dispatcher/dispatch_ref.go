package http_server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func handleRefProfilePictures(kernel interfaces.Kernel, w http.ResponseWriter, r *http.Request, username string) {
	data, err := kernel.BlobManager().ReadBlob("profile_picture", username, "", "")
	if err != nil {
		http.Redirect(w, r, "https://static.hereus.net/protocols/profile_photos/_default.png", http.StatusFound)
		return
	}
	mimeType := http.DetectContentType(data)
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func dispatchRef(kernel interfaces.Kernel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint, _ := strings.CutPrefix(r.URL.Path, "/ref/")
		endpoint = strings.TrimSuffix(endpoint, "/")
		endpoint = strings.ReplaceAll(endpoint, "//", "/")
		if len(endpoint) == 0 {
			http.Error(w, "Endpoint not specified", http.StatusBadRequest)
			return
		}
		fmt.Printf("(blob) %s [%s] %s\n", r.RemoteAddr, time.Now().Format("2006-01-02 15:04:05"), endpoint)
		endpointPieces := strings.Split(endpoint, "/")
		if len(endpointPieces) < 1 {
			http.Error(w, "Invalid endpoint", http.StatusBadRequest)
			return
		}
		switch endpointPieces[0] {
		case "accounts":
			switch endpointPieces[1] {
			case "profile_pictures":
				handleRefProfilePictures(kernel, w, r, endpointPieces[2])
			}
		default:
			http.Error(w, fmt.Sprintf("Reference not found: %s", endpoint), http.StatusNotFound)
		}
	}
}
