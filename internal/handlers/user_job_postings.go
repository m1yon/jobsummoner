package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/lmittmann/tint"
	"github.com/m1yon/jobsummoner/internal/database"
)

func (cfg *handlersConfig) patchUserJobPostingsHandler(w http.ResponseWriter, r *http.Request) {
	jobPostingID := r.PathValue("jobPostingID")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	currentStatus := r.FormValue("currentStatus")
	status := r.FormValue("status")

	slog.Info("", slog.String("currentStatus", currentStatus), slog.String("status", status))

	CurrentStatus, err := strconv.Atoi(currentStatus)

	if err != nil {
		slog.Error("failed to convert currentStatus", tint.Err(err))
		CurrentStatus = 1
	}

	Status, err := strconv.Atoi(status)

	if err != nil {
		slog.Error("failed to convert status", tint.Err(err))
		Status = 1
	}

	err = cfg.DB.UpdateUserJobPostingStatus(r.Context(), database.UpdateUserJobPostingStatusParams{UserID: 1, JobPostingID: jobPostingID, Status: int64(Status)})

	if err != nil {
		slog.Error("failed to hide user job posting", slog.String("id", jobPostingID), tint.Err(err))
	}

	JobPostings, err := cfg.DB.GetUserJobPostingsByStatus(r.Context(), database.GetUserJobPostingsByStatusParams{UserID: 1, Status: int64(CurrentStatus)})

	if err != nil {
		slog.Error("failed to query job postings", tint.Err(err))
	}

	formattedJobPostings := make([]FormattedJobPosting, 0, len(JobPostings))

	for _, jobPosting := range JobPostings {
		formattedJobPostings = append(formattedJobPostings, FormattedJobPosting{GetUserJobPostingsByStatusRow: jobPosting, TimeAgo: timeAgo(jobPosting.LastPosted)})
	}

	component := userJobPostingsTemplate(formattedJobPostings, CurrentStatus)
	component.Render(context.Background(), w)

	if err != nil {
		slog.Error("could not execute template", tint.Err(err))
	}
}

func (cfg *handlersConfig) getUserJobPostingsHandler(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")

	Status, err := strconv.Atoi(statusFilter)

	if err != nil {
		Status = 1
	}

	JobPostings, err := cfg.DB.GetUserJobPostingsByStatus(r.Context(), database.GetUserJobPostingsByStatusParams{UserID: 1, Status: int64(Status)})

	if err != nil {
		slog.Error("failed to query job postings", tint.Err(err))
	}

	formattedJobPostings := make([]FormattedJobPosting, 0, len(JobPostings))

	for _, jobPosting := range JobPostings {
		formattedJobPostings = append(formattedJobPostings, FormattedJobPosting{GetUserJobPostingsByStatusRow: jobPosting, TimeAgo: timeAgo(jobPosting.LastPosted)})
	}

	component := userJobPostingsTemplate(formattedJobPostings, Status)
	component.Render(context.Background(), w)

	if err != nil {
		slog.Error("could not execute template", tint.Err(err))
	}
}
