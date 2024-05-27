package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGETHomepage(t *testing.T) {
	t.Run("renders the available positions correctly", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server := DefaultHomepageServer{}
		server.Get(response, request)

		assert.Equal(t, response.Code, 200)
		assert.Equal(t, response.Body.String(), "<div><p>Software Engineer,Manager,</p></div>")
	})
}
