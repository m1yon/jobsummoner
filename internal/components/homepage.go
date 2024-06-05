package components

import (
	"github.com/m1yon/jobsummoner"
)

func NewHomepageViewModel(jobs []jobsummoner.Job) jobsummoner.HomepageViewModel {
	m := jobsummoner.HomepageViewModel{
		Jobs: jobs,
	}

	return m
}
