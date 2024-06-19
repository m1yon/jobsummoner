package main

import (
	"log/slog"
	"net/http"

	"github.com/m1yon/jobsummoner"
)

func newHttpClient(logger *slog.Logger, config *config) *http.Client {
	httpClient, err := jobsummoner.NewHttpProxyClientFromConfig(config.proxyConfig)

	if err != nil {
		logger.Warn("proxy server disabled", "reason", err.Error())
	} else {
		logger.Info("proxy server enabled")
	}

	return httpClient
}
