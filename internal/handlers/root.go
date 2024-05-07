package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/lmittmann/tint"
	"github.com/m1yon/jobsummoner/internal/database"
)

func (cfg *handlersConfig) rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	JobPostings, err := cfg.DB.GetUserJobPostings(ctx, 1)

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

	lastScrapedDate, err := cfg.DB.GetLastScrapedDate(ctx, 1)

	if err != nil {
		slog.Error("failed to get user's last scraped date", tint.Err(err))
	}

	LastScrapedTime := "N/A"

	if lastScrapedDate.Valid {
		LastScrapedTime = timeAgo(lastScrapedDate.Time)
	}

	homepage := struct {
		JobPostings     []FormattedJobPosting
		LastScrapedTime string
	}{
		formattedJobPostings,
		LastScrapedTime,
	}

	cfg.Renderer.Render(w, "root", homepage)

	if err != nil {
		slog.Error("could not execute template", tint.Err(err))
	}
}

func timeAgo(from time.Time) string {
	now := time.Now()
	diff := now.Sub(from)

	if diff < time.Minute {
		if int(diff.Seconds()) == 1 {
			return fmt.Sprintf("%d second ago", int(diff.Seconds()))
		}
		return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
	} else if diff < time.Hour {
		if int(diff.Minutes()) == 1 {
			return fmt.Sprintf("%d minute ago", int(diff.Minutes()))
		}
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	} else if diff < time.Hour*24 {
		if int(diff.Hours()) == 1 {
			return fmt.Sprintf("%d hour ago", int(diff.Hours()))
		}
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	} else if diff < time.Hour*24*30 {
		days := diff / (time.Hour * 24)
		if days == 1 {
			return fmt.Sprintf("%d day ago", int(days))
		}
		return fmt.Sprintf("%d days ago", int(days))
	} else if diff < time.Hour*24*365 {
		months := diff / (time.Hour * 24 * 30)
		if months == 1 {
			return fmt.Sprintf("%d month ago", int(months))
		}
		return fmt.Sprintf("%d months ago", int(months))
	}
	years := diff / (time.Hour * 24 * 365)
	if years == 1 {
		return fmt.Sprintf("%d year ago", int(years))
	}
	return fmt.Sprintf("%d years ago", int(years))
}
