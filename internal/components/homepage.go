package components

import (
	"github.com/m1yon/jobsummoner"
)

func NewHomepageViewModel(jobs []jobsummoner.Job) (m jobsummoner.HomepageViewModel) {
	for _, position := range jobs {
		m.Text += position.Name + ","
	}

	return m
}
