package main

type HomepageViewModel struct {
	text string
}

func NewHomepageViewModel(jobs []Job) (m HomepageViewModel) {
	for _, position := range jobs {
		m.text += position.name + ","
	}

	return m
}
