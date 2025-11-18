package rpc_org_theprotocols_application

import (
	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func SaveDataRpc() *helpers.RpcFunctionHandlerWithArguments {
	return &helpers.RpcFunctionHandlerWithArguments{
		ReqFactory: func() interface{} {
			return &SaveDataArguments{}
		},
		Handler: func(session interfaces.Session, req interface{}) (interface{}, error) {
			return nil, SaveData(session, req.(*SaveDataArguments))
		},
		Permissions: []string{},
	}
}

type SaveDataArguments struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

func SaveData(session interfaces.Session, req *SaveDataArguments) error {
	app := session.GetApp()
	return app.SaveData(req.Name, req.Data)
}
