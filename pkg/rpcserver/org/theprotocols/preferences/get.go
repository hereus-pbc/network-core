package org_theprotocols_preferences

import (
	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func GetRpc() *helpers.RpcFunctionHandlerWithArguments {
	return &helpers.RpcFunctionHandlerWithArguments{
		ReqFactory: func() interface{} {
			return &GetArguments{}
		},
		Handler: func(session interfaces.Session, req interface{}) (interface{}, error) {
			return Get(session, req.(*GetArguments))
		},
		Permissions: []string{},
	}
}

type GetArguments struct {
	PackageName string `json:"packageName"`
}

func Get(session interfaces.Session, req *GetArguments) (map[string]interface{}, error) {
	app, err := session.GetUser().App(req.PackageName)
	if err != nil {
		return nil, err
	}
	return app.GetPreferences()
}
