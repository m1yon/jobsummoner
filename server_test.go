package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
)

func TestGETHomepage(t *testing.T) {
	t.Run("renders the available positions correctly", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server := NewDefaultServer()
		server.ServerHTTP(response, request)

		assert.Equal(t, response.Code, 200)
		assert.Equal(t, response.Body.String(), "<div><p>Software Engineer,Manager,</p></div>")
	})

	t.Run("handles a template rendering failure", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server := NewDefaultServer()
		server.Render = func(component templ.Component, ctx context.Context, w io.Writer) error {
			return errors.New("could not render template")
		}
		server.ServerHTTP(response, request)

		assert.Equal(t, 500, response.Code)
		assert.Equal(t, "", response.Body.String())
	})
}
