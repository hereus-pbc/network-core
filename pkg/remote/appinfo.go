package remote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hereus-pbc/network-core/pkg/types"
)

//goland:noinspection ALL
func FetchAppInfo(packageName string) *types.AppInfo {
	packageNamePieces := strings.Split(packageName, ".")
	domainPieces := make([]string, len(packageNamePieces))
	for i := len(packageNamePieces) - 1; i >= 0; i-- {
		domainPieces[len(packageNamePieces)-1-i] = packageNamePieces[i]
	}
	domain := strings.Join(domainPieces, ".")
	resp, err := http.Get(fmt.Sprintf("https://%s/.well-known/app_info.json", domain))
	if err != nil {
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close app info response body:", err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	var appInfo types.AppInfo
	err = json.NewDecoder(resp.Body).Decode(&appInfo)
	if err != nil {
		return nil
	}
	return &appInfo
}
