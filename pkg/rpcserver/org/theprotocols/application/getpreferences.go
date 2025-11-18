package rpc_org_theprotocols_application

import (
	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func GetPreferencesRpc() *helpers.RpcFunctionHandlerNoArguments {
	return &helpers.RpcFunctionHandlerNoArguments{
		Handler: func(session interfaces.Session) (interface{}, error) {
			return GetPreferences(session)
		},
		Permissions: []string{},
	}
}

func GetPreferences(session interfaces.Session) (map[string]interface{}, error) {
	app := session.GetApp()
	return app.GetPreferences()
}
