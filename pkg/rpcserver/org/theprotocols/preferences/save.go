package org_theprotocols_preferences

import (
	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func SaveRpc() *helpers.RpcFunctionHandlerWithArguments {
	return &helpers.RpcFunctionHandlerWithArguments{
		ReqFactory: func() interface{} {
			return &SavePreferencesArguments{}
		},
		Handler: func(session interfaces.Session, req interface{}) (interface{}, error) {
			return nil, Save(session, req.(*SavePreferencesArguments))
		},
		Permissions: []string{},
	}
}

type SavePreferencesArguments struct {
	PackageName string                 `json:"packageName"`
	Preferences map[string]interface{} `json:"preferences"`
}

func Save(session interfaces.Session, req *SavePreferencesArguments) error {
	app, err := session.GetUser().App(req.PackageName)
	if err != nil {
		return err
	}
	return app.SavePreferences(req.Preferences)
}
