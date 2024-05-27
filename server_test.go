package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGETHomepage(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	HomepageServer(response, request)

	assert.Equal(t, response.Body.String(), "Software Engineer,Manager,")
}
