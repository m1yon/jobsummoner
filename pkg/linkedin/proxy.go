package linkedin

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type ProxyConfig struct {
	Hostname string
	Port     string
	Username string
	Password string
}

func BuildHttpProxyURL(config ProxyConfig) (*url.URL, error) {
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