package rpc_org_theprotocols_application

import (
	"github.com/hereus-pbc/network-core/pkg/httpserver/helpers"
	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

func GetDataRpc() *helpers.RpcFunctionHandlerWithArguments {
	return &helpers.RpcFunctionHandlerWithArguments{
		ReqFactory: func() interface{} {
			return &GetDataArguments{}
		},
		Handler: func(session interfaces.Session, req interface{}) (interface{}, error) {
			return GetData(session, req.(*GetDataArguments))
		},
		Permissions: []string{},
	}
}

type GetDataArguments struct {
	Name string `json:"name"`
}

func GetData(session interfaces.Session, req *GetDataArguments) ([]byte, error) {
	app := session.GetApp()
	return app.GetData(req.Name)
}
