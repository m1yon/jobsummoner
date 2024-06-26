package components

import "github.com/m1yon/jobsummoner/internal/models"

type HomepageJobModel struct {
	models.Job
	LastPostedText string
}

type HomepageViewModel struct {
	Jobs            []HomepageJobModel
	Flash           string
	IsAuthenticated bool
}
