package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/m1yon/jobsummoner/pkg/linkedin"
)

func newHttpClient(logger *slog.Logger) *http.Client {
	proxyConfig := linkedin.ProxyConfig{
		Hostname: os.Getenv("PROXY_HOSTNAME"),
		Port:     os.Getenv("PROXY_PORT"),
		Username: os.Getenv("PROXY_USERNAME"),
		Password: os.Getenv("PROXY_PASSWORD"),
	}

	httpClient, err := linkedin.NewHttpProxyClientFromConfig(proxyConfig)

	if err != nil {
		logger.Warn("proxy server disabled", "reason", err.Error())
	} else {
		logger.Info("proxy server enabled")
	}

	return httpClient
}
