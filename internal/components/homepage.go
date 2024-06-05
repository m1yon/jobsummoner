package components

import (
	"fmt"
	"time"

	"github.com/m1yon/jobsummoner"
)

func NewHomepageViewModel(jobs []jobsummoner.Job) jobsummoner.HomepageViewModel {
	jobModels := make([]jobsummoner.HomepageJobModel, 0, len(jobs))

	for _, job := range jobs {
		jobModels = append(jobModels, jobsummoner.HomepageJobModel{
			Job:            job,
			LastPostedText: timeAgo(job.LastPosted),
		})
	}

	m := jobsummoner.HomepageViewModel{
		Jobs: jobModels,
	}

	return m
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
