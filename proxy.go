package jobsummoner

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	ErrInvalidProxyConfig = "invalid proxy config provided"
	ErrBuildingProxyURL   = "failed to build http proxy url"
	ErrProxySetup         = "problem setting up http proxy client"
)

type ProxyConfig struct {
	Hostname string
	Port     string
	Username string
	Password string
}

func BuildHttpProxyURL(config ProxyConfig) (*url.URL, error) {
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

func NewHttpProxyClient(proxyURL *url.URL) (*http.Client, error) {
	proxy := http.ProxyURL(proxyURL)
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport}

	return client, nil
}

func NewHttpProxyClientFromConfig(config ProxyConfig) (*http.Client, error) {
	proxyURL, err := BuildHttpProxyURL(config)

	if err != nil {
		return &http.Client{}, errors.Wrap(err, ErrBuildingProxyURL)
	}

	client, err := NewHttpProxyClient(proxyURL)

	if err != nil {
		return &http.Client{}, errors.Wrap(err, ErrProxySetup)
	}

	return client, nil
}
