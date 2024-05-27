package main

import (
	"context"
	"io"
)

type HomepageViewModel struct {
	text string
}

func NewHomepageViewModel(jobPostings []JobPosting) (m HomepageViewModel) {
	for _, position := range jobPostings {
		m.text += position.name + ","
	}

	return m
}

func RenderHomepage(m HomepageViewModel, w io.Writer) {
	component := homepage(m)
	component.Render(context.Background(), w)
}
