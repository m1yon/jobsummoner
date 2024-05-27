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
}

func (h *DefaultHomepageServer) Get(w http.ResponseWriter, r *http.Request) {
	jobPostings := GetJobPostings()
	m := NewHomepageViewModel(jobPostings)

	component := homepage(m)
	h.Render(component, context.Background(), w)
}

func (h *DefaultHomepageServer) Render(component templ.Component, ctx context.Context, w io.Writer) error {
	return component.Render(ctx, w)
}
