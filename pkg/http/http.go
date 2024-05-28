package http

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/components"
)

type Server interface {
	Render(component templ.Component, ctx context.Context, w io.Writer) error
	ListenAndServe(addr string)
}

type DefaultServer struct {
	Render func(component templ.Component, ctx context.Context, w io.Writer) error
	jobsummoner.JobService
}

func NewDefaultServer() *DefaultServer {
	return &DefaultServer{
		Render:     Render,
		JobService: jobsummoner.NewDefaultJobService(),
	}
}

func (h *DefaultServer) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	h.GetHomepage(w, r)
}

func (h *DefaultServer) GetHomepage(w http.ResponseWriter, r *http.Request) {
	jobs := h.JobService.GetJobs()

	m := jobsummoner.NewHomepageViewModel(jobs)
	component := components.Homepage(m)
	err := h.Render(component, context.Background(), w)

	if err != nil {
		w.WriteHeader(500)
	}
}

func (h *DefaultServer) ListenAndServe(addr string) {
	handler := http.HandlerFunc(h.ServerHTTP)
	log.Fatal(http.ListenAndServe(addr, handler))
}

func Render(component templ.Component, ctx context.Context, w io.Writer) error {
	return component.Render(ctx, w)
}
