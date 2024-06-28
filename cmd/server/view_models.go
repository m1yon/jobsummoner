package main

import (
	"net/http"

	"github.com/m1yon/jobsummoner/internal/components"
	"github.com/m1yon/jobsummoner/internal/models"
)

func (app *application) NewHomepageViewModel(r *http.Request, jobs []models.Job) components.HomepageViewModel {
	flash := app.sessionManager.PopString(r.Context(), "flash")
	jobModels := make([]components.HomepageJobModel, 0, len(jobs))

	for _, job := range jobs {
		jobModels = append(jobModels, components.HomepageJobModel{
			Job:            job,
			LastPostedText: timeAgo(job.LastPosted),
		})
	}

	m := components.HomepageViewModel{
		Jobs:            jobModels,
		Flash:           flash,
		IsAuthenticated: app.isAuthenticated(r),
	}

	return m
}
