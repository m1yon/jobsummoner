package main

import (
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
)

type Server interface {
	Render(component templ.Component, ctx context.Context, w io.Writer) error
}

type DefaultServer struct {
	Render func(component templ.Component, ctx context.Context, w io.Writer) error
}

func NewDefaultServer() *DefaultServer {
	return &DefaultServer{
		Render: Render,
	}
}

func (h *DefaultServer) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	h.GetHomepage(w, r)
}

func (h *DefaultServer) GetHomepage(w http.ResponseWriter, r *http.Request) {
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
