package http

import (
	"net/http"

	"github.com/m1yon/jobsummoner"
)

func (s *Server) NewHomepageViewModel(r *http.Request, jobs []jobsummoner.Job) jobsummoner.HomepageViewModel {
	flash := s.sessionManager.PopString(r.Context(), "flash")
	jobModels := make([]jobsummoner.HomepageJobModel, 0, len(jobs))

	for _, job := range jobs {
		jobModels = append(jobModels, jobsummoner.HomepageJobModel{
			Job:            job,
			LastPostedText: timeAgo(job.LastPosted),
		})
	}

	m := jobsummoner.HomepageViewModel{
		Jobs:            jobModels,
		Flash:           flash,
		IsAuthenticated: s.isAuthenticated(r),
	}

	return m
}
