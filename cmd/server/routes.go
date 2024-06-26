package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	dynamic := alice.New(s.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(s.getHomepageHandler))

	mux.Handle("GET /user/signup", dynamic.ThenFunc(s.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(s.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(s.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(s.userLoginPost))
	mux.Handle("POST /user/logout", dynamic.ThenFunc(s.userLogoutPost))

	standard := alice.New(s.logRequest, s.recoverPanic, commonHeaders)
	return standard.Then(mux)
}
