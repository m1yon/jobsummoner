package http

import (
	"net/http"

	"github.com/justinas/alice"
)

func (server *DefaultServer) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", server.getHomepageHandler)

	standard := alice.New(server.logRequest, server.recoverPanic, commonHeaders)

	return standard.Then(mux)
}
