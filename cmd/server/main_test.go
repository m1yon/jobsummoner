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
	driver, err := web.NewWebDriver(port)

	if err != nil {
		t.Fatal("failed creating web driver:", err.Error())
	}

	tests.StartDockerServer(t, port)
	specifications.AuthSpecification(t, driver)
}
