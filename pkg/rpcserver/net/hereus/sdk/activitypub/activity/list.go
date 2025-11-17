package rpc_net_hereus_sdk_activitypub_activity

import (
	"fmt"

	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func ListRpc() *helpers.RpcFunctionHandlerWithArguments {
	return &helpers.RpcFunctionHandlerWithArguments{
		ReqFactory: func() interface{} { return &ListArguments{} },
		Handler: func(session interfaces.Session, req interface{}) (interface{}, error) {
			return List(session, req.(*ListArguments))
		},
		Permissions: []string{"net.hereus.sdk.permissions.activitypub"},
	}
}

type ListArguments struct {
	Filters         map[string]interface{} `json:"filters"`
	PackageSpecific bool                   `json:"packageSpecific"`
	Max             int64                  `json:"max"`
	LatestFirst     bool                   `json:"latestFirst"`
}

func List(session interfaces.Session, req *ListArguments) (*[]ActivityType, error) {
	packageName := (func() string {
		if req.PackageSpecific {
			return session.GetSession().PackageName
		}
		return ""
	})()
	results, err := session.GetKernel().ActivityPubDB().ListNotes(session.GetUser().GetActorUrl(), req.Filters, packageName, req.Max, req.LatestFirst)
	if err != nil {
		return nil, err
	}
	var resultsSdk []ActivityType
	for i := range results {
		results[i].CustomUrl = ""
		sdkObject, err := NoteToSdkObject(session.GetKernel(), results[i])
		if err != nil {
			fmt.Println("Failed to convert note to SDK object:", results[i])
		} else {
			resultsSdk = append(resultsSdk, *sdkObject)
		}
	}
	if len(resultsSdk) == 0 {
		resultsSdk = []ActivityType{}
	}
	return &resultsSdk, nil
}
