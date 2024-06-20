package http

import (
	"context"
	"net/http"

	"github.com/m1yon/jobsummoner/internal/components"
)

func (s *Server) getHomepageHandler(w http.ResponseWriter, r *http.Request) {
	jobs, err := s.jobService.GetJobs(r.Context())

	if err != nil {
		s.serverError(w, r, err)
	}

	m := components.NewHomepageViewModel(jobs)
	component := components.Homepage(m)
	err = s.Render(component, context.Background(), w)

	if err != nil {
		s.serverError(w, r, err)
	}
}
