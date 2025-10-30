package http_server

import (
	"fmt"
	"net/http"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
)

type Route struct {
	Handler http.HandlerFunc
	Enabled bool
}

type HttpServer struct {
	kernel       interfaces.Kernel
	port         int
	mux          *http.ServeMux
	customRoutes map[string]Route
}

//goland:noinspection GoUnusedExportedFunction
func CreateServer(kernel interfaces.Kernel, port int) *HttpServer {
	s := &HttpServer{
		kernel:       kernel,
		port:         port,
		mux:          http.NewServeMux(),
		customRoutes: make(map[string]Route),
	}
	s.mux.HandleFunc("/.well-known/app_info.json", handleAppInfo(kernel))
	s.mux.HandleFunc("/theprotocols/", dispatchRpcEndpoints(kernel))
	s.mux.HandleFunc("/ref/", dispatchRef(kernel))
	s.mux.HandleFunc("/activitypub/", dispatchActivityPub(kernel))
	s.mux.HandleFunc("/.well-known/webfinger", handleWebFinger(kernel))
	return s
}

func (s *HttpServer) HandleFunc(path string, handler http.HandlerFunc) {
	if _, exists := map[string]bool{
		"/.well-known/app_info.json": true,
		"/theprotocols/":             true,
		"/ref/":                      true,
		"/activitypub/":              true,
		"/.well-known/webfinger":     true,
	}[path]; exists {
		panic(fmt.Sprintf("Cannot override built-in route: %s", path))
	}
	if _, exists := s.customRoutes[path]; exists {
		panic(fmt.Sprintf("Route already exists: %s", path))
	}
	s.customRoutes[path] = Route{
		Handler: handler,
		Enabled: true,
	}
	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		route, exists := s.customRoutes[path]
		if !exists {
			http.Error(w, "Route not found", http.StatusNotFound)
		}
		if !route.Enabled {
			http.Error(w, "Route disabled", http.StatusNotFound)
			return
		}
		route.Handler(w, r)
	})
}

func (s *HttpServer) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}
