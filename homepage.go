package jobsummoner

import "github.com/m1yon/jobsummoner/internal/models"

func NewHomepageViewModel(jobs []Job) (m models.HomepageViewModel) {
	for _, position := range jobs {
		m.Text += position.name + ","
	}

	return m
}
