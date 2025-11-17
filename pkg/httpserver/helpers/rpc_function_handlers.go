package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

type RpcFunctionHandler interface {
	Handle(kernel interfaces.Kernel, r *http.Request, w http.ResponseWriter)
}

type RpcFunctionMapping map[string]RpcFunctionHandler

func AutoHandleRpcFunction(mapped *RpcFunctionMapping, kernel interfaces.Kernel, r *http.Request, w http.ResponseWriter, functionName string) error {
	if handler, exists := (*mapped)[functionName]; exists {
		handler.Handle(kernel, r, w)
		return nil
	}
	return fmt.Errorf("unknown endpoint: %s", functionName)
}

type RpcFunctionHandlerWithArguments struct {
	ReqFactory  func() interface{}
	Handler     func(interfaces.Session, interface{}) (interface{}, error)
	Permissions []string
}

func (fh *RpcFunctionHandlerWithArguments) Handle(kernel interfaces.Kernel, r *http.Request, w http.ResponseWriter) {
	req := fh.ReqFactory()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	kernel.SessionWrapper(w, r, fh.Permissions, func(session interfaces.Session) {
		output, err := fh.Handler(session, req)
		ConvertRpcResponseToHttpResponse(output, err, w)
	})
}

type RpcFunctionHandlerNoArguments struct {
	Handler     func(interfaces.Session) (interface{}, error)
	Permissions []string
}

func (fh *RpcFunctionHandlerNoArguments) Handle(kernel interfaces.Kernel, r *http.Request, w http.ResponseWriter) {
	kernel.SessionWrapper(w, r, fh.Permissions, func(session interfaces.Session) {
		output, err := fh.Handler(session)
		ConvertRpcResponseToHttpResponse(output, err, w)
	})
}

type RpcFunctionHandlerNoArgumentsNoSession struct {
	Handler func(interfaces.Kernel) (interface{}, error)
}

func (fh *RpcFunctionHandlerNoArgumentsNoSession) Handle(kernel interfaces.Kernel, r *http.Request, w http.ResponseWriter) {
	output, err := fh.Handler(kernel)
	ConvertRpcResponseToHttpResponse(output, err, w)
}
