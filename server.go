package main

import (
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
)

type HomepageServer interface {
	Render(component templ.Component, ctx context.Context, w io.Writer) error
}

type DefaultHomepageServer struct {
	Render func(component templ.Component, ctx context.Context, w io.Writer) error
}

func NewDefaultHomepageServer() *DefaultHomepageServer {
	return &DefaultHomepageServer{
		Render: Render,
	}
}

func (h *DefaultHomepageServer) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	h.Get(w, r)
}

func (h *DefaultHomepageServer) Get(w http.ResponseWriter, r *http.Request) {
	jobPostings := GetJobPostings()
	m := NewHomepageViewModel(jobPostings)

	component := homepage(m)
	err := h.Render(component, context.Background(), w)

	if err != nil {
		w.WriteHeader(500)
	}
}

func Render(component templ.Component, ctx context.Context, w io.Writer) error {
	return component.Render(ctx, w)
}
