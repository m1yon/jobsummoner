package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/m1yon/jobsummoner/internal/components"
)

func (h *DefaultServer) getHomepageHandler(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.JobService.GetJobs(r.Context())

	if err != nil {
		h.logger.Error("problem getting jobs", slog.String("err", err.Error()))
		w.WriteHeader(500)
	}

	m := components.NewHomepageViewModel(jobs)
	component := components.Homepage(m)
	err = h.Render(component, context.Background(), w)

	if err != nil {
		w.WriteHeader(500)
	}
}
