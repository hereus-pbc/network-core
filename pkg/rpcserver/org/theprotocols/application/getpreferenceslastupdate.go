package rpc_org_theprotocols_application

import (
	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func GetPreferencesLastUpdateRpc() *helpers.RpcFunctionHandlerNoArguments {
	return &helpers.RpcFunctionHandlerNoArguments{
		Handler: func(session interfaces.Session) (interface{}, error) {
			return GetPreferencesLastUpdate(session), nil
		},
		Permissions: []string{},
	}
}

func GetPreferencesLastUpdate(session interfaces.Session) int64 {
	app := session.GetApp()
	return app.GetPreferencesLastUpdate()
}
