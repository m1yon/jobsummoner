package jobsummoner

import (
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/m1yon/jobsummoner/internal/components"
)

type Server interface {
	Render(component templ.Component, ctx context.Context, w io.Writer) error
}

type DefaultServer struct {
	Render func(component templ.Component, ctx context.Context, w io.Writer) error
	JobService
}

func NewDefaultServer() *DefaultServer {
	return &DefaultServer{
		Render:     Render,
		JobService: NewDefaultJobService(),
	}
}

func (h *DefaultServer) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	h.GetHomepage(w, r)
}

func (h *DefaultServer) GetHomepage(w http.ResponseWriter, r *http.Request) {
	jobs := h.JobService.GetJobs()

	m := NewHomepageViewModel(jobs)
	component := components.Homepage(m)
	err := h.Render(component, context.Background(), w)

	if err != nil {
		w.WriteHeader(500)
	}
}

func Render(component templ.Component, ctx context.Context, w io.Writer) error {
	return component.Render(ctx, w)
}
