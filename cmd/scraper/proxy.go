package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	ErrInvalidProxyConfig = "invalid proxy config provided"
	ErrBuildingProxyURL   = "failed to build http proxy url"
	ErrProxySetup         = "problem setting up http proxy client"
)

type proxyConfig struct {
	Hostname string
	Port     string
	Username string
	Password string
}

func buildHttpProxyURL(config proxyConfig) (*url.URL, error) {
	if config.Hostname == "" || config.Port == "" {
		return &url.URL{}, errors.New(ErrInvalidProxyConfig)
	}

	urlString := fmt.Sprintf("http://%v:%v@%v:%v", config.Username, config.Password, config.Hostname, config.Port)
	parsedURL, err := url.Parse(urlString)

	if err != nil {
		return &url.URL{}, errors.Wrap(err, "unable to parse proxyURL")
	}

	return parsedURL, nil
}

func newHttpProxyClient(proxyURL *url.URL) (*http.Client, error) {
	proxy := http.ProxyURL(proxyURL)
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport}

	return client, nil
}

func newHttpProxyClientFromConfig(config proxyConfig) (*http.Client, error) {
	proxyURL, err := buildHttpProxyURL(config)

	if err != nil {
		return &http.Client{}, errors.Wrap(err, ErrBuildingProxyURL)
	}

	client, err := newHttpProxyClient(proxyURL)

	if err != nil {
		return &http.Client{}, errors.Wrap(err, ErrProxySetup)
	}

	return client, nil
}

func newHttpClient(logger *slog.Logger, config *config) *http.Client {
	httpClient, err := newHttpProxyClientFromConfig(config.proxyConfig)

	if err != nil {
		logger.Warn("proxy server disabled", "reason", err.Error())
	} else {
		logger.Info("proxy server enabled")
	}

	return httpClient
}
