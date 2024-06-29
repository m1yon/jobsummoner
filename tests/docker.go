package tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	startupTimeout = 1 * time.Minute
	dockerFilePath = "docker/server/Dockerfile"
)

type StdoutLogConsumer struct{}

func (lc *StdoutLogConsumer) Accept(l testcontainers.Log) {
	fmt.Print(string(l.Content))
}

func StartDockerServer(
	t testing.TB,
	port string,
) (string, error) {
	t.Helper()

	g := StdoutLogConsumer{}

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		FromDockerfile: newTCDockerfile(),
		ExposedPorts:   []string{port},
		WaitingFor:     wait.ForLog("server started").WithStartupTimeout(startupTimeout),
		Env: map[string]string{
			"LOCAL_DB": "true",
		},
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Opts:      []testcontainers.LogProductionOption{testcontainers.WithLogProductionTimeout(10 * time.Second)},
			Consumers: []testcontainers.LogConsumer{&g},
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	baseURL, err := container.Endpoint(ctx, "")
	assert.NoError(t, err)

	res, err := http.Get("http://" + baseURL)
	assert.NoError(t, err)

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		body = bytes.TrimSpace(body)

		return "", fmt.Errorf("received failing status code: %v\n%v", res.StatusCode, string(body))
	}

	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})

	return baseURL, nil
}

func newTCDockerfile() testcontainers.FromDockerfile {
	return testcontainers.FromDockerfile{
		Context:       ".",
		Dockerfile:    dockerFilePath,
		PrintBuildLog: true,
	}
}
