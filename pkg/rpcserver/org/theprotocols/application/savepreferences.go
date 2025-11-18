package rpc_org_theprotocols_application

import (
	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func SavePreferencesRpc() *helpers.RpcFunctionHandlerWithArguments {
	return &helpers.RpcFunctionHandlerWithArguments{
		ReqFactory: func() interface{} {
			return &SavePreferencesArguments{}
		},
		Handler: func(session interfaces.Session, req interface{}) (interface{}, error) {
			return nil, SavePreferences(session, req.(*SavePreferencesArguments))
		},
		Permissions: []string{},
	}
}

type SavePreferencesArguments struct {
	Preferences map[string]interface{} `json:"preferences"`
}

func SavePreferences(session interfaces.Session, req *SavePreferencesArguments) error {
	app := session.GetApp()
	return app.SavePreferences(req.Preferences)
}
