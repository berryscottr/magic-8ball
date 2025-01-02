package healthz

import "net/http"

type HealthCheckServer struct {
	server *http.Server
}

type Methods interface {
	Start()
	Close()
}