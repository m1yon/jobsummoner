package handlers

import (
	"log/slog"
	"net/http"

	"github.com/lmittmann/tint"
	"github.com/m1yon/jobsummoner/internal/database"
)

func (cfg *handlersConfig) putUserJobPostingsHandler(w http.ResponseWriter, r *http.Request) {
	jobPostingID := r.PathValue("jobPostingID")

	err := cfg.DB.HideUserJobPosting(r.Context(), database.HideUserJobPostingParams{UserID: 1, JobPostingID: jobPostingID})

	if err != nil {
		slog.Error("failed to hide user job posting", slog.String("id", jobPostingID), tint.Err(err))
	}

	JobPostings, err := cfg.DB.GetUserJobPostings(r.Context(), 1)

	if err != nil {
		slog.Error("failed to query job postings", tint.Err(err))
	}

	type FormattedJobPosting struct {
		database.GetUserJobPostingsRow
		TimeAgo string
	}

	formattedJobPostings := make([]FormattedJobPosting, 0, len(JobPostings))

	for _, jobPosting := range JobPostings {
		formattedJobPostings = append(formattedJobPostings, FormattedJobPosting{GetUserJobPostingsRow: jobPosting, TimeAgo: timeAgo(jobPosting.LastPosted)})
	}

	FormattedJobPostings := struct {
		JobPostings []FormattedJobPosting
	}{
		formattedJobPostings,
	}

	err = cfg.Renderer.Render(w, "user_job_postings", FormattedJobPostings)

	if err != nil {
		slog.Error("could not execute template", tint.Err(err))
	}
}
