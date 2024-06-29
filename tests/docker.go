package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	startupTimeout = 5 * time.Second
	dockerFilePath = "docker/server/Dockerfile"
)

func StartDockerServer(
	t testing.TB,
	port string,
) {
	t.Helper()

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		FromDockerfile: newTCDockerfile(),
		ExposedPorts:   []string{fmt.Sprintf("%s:%s", port, port)},
		WaitingFor:     wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(startupTimeout),
		Env: map[string]string{
			"LOCAL_DB": "true",
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})
}

func newTCDockerfile() testcontainers.FromDockerfile {
	return testcontainers.FromDockerfile{
		Context:       ".",
		Dockerfile:    dockerFilePath,
		PrintBuildLog: true,
	}
}
