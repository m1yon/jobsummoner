package main

import (
	"testing"

	"github.com/m1yon/jobsummoner/tests"
	"github.com/m1yon/jobsummoner/tests/adapters/web"
	"github.com/m1yon/jobsummoner/tests/specifications"
)

func TestServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	port := "3000"

	baseURL, err := tests.StartDockerServer(t, port)

	if err != nil {
		t.Fatal(err)
	}

	driver, err := web.NewWebDriver(baseURL)

	if err != nil {
		t.Fatal("failed creating web driver:", err.Error())
	}

	specifications.AuthSpecification(t, driver)
}
