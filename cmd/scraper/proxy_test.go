package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpProxyBuilder(t *testing.T) {
	t.Run("builds a valid URL", func(t *testing.T) {
		fakeProxyConfig := proxyConfig{
			Hostname: "crazyproxies.com",
			Port:     "44443",
			Username: "michael",
			Password: "hunter12",
		}

		url, err := buildHttpProxyURL(fakeProxyConfig)
		assert.NoError(t, err)
		assert.Equal(t, "http://michael:hunter12@crazyproxies.com:44443", url.String())
	})

	t.Run("returns an error on empty values", func(t *testing.T) {
		fakeProxyConfig := proxyConfig{
			Hostname: "",
			Port:     "",
			Username: "",
			Password: "",
		}

		_, err := buildHttpProxyURL(fakeProxyConfig)
		assert.ErrorContains(t, err, ErrInvalidProxyConfig)
	})
}

func TestHttpProxyClient(t *testing.T) {
	t.Run("initializes proxy correctly", func(t *testing.T) {
		proxyCalled := false
		proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxyCalled = true
			w.WriteHeader(http.StatusOK)
		}))
		defer proxyServer.Close()

		proxyServerURL, _ := url.Parse(proxyServer.URL)
		client, _ := newHttpProxyClient(proxyServerURL)

		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
		_, _ = client.Transport.RoundTrip(req)

		assert.Equal(t, true, proxyCalled, "proxy was not called")
	})
}

func TestNewHttpProxyClientFromConfig(t *testing.T) {
	t.Run("initializes proxy client correctly", func(t *testing.T) {
		proxyCalled := false
		proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxyCalled = true
			w.WriteHeader(http.StatusOK)
		}))
		defer proxyServer.Close()

		// split into hostname/port
		proxyInfo := strings.Split(strings.ReplaceAll(proxyServer.URL, "http://", ""), ":")

		client, _ := newHttpProxyClientFromConfig(proxyConfig{
			Hostname: proxyInfo[0],
			Port:     proxyInfo[1],
		})

		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
		_, _ = client.Transport.RoundTrip(req)

		assert.Equal(t, true, proxyCalled, "proxy was not called")
	})

	t.Run("returns an error for bad config", func(t *testing.T) {
		proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer proxyServer.Close()

		// split into hostname/port
		proxyInfo := strings.Split(strings.ReplaceAll(proxyServer.URL, "http://", ""), ":")

		_, err := newHttpProxyClientFromConfig(proxyConfig{
			Hostname: "",
			Port:     proxyInfo[1],
		})
		assert.ErrorContains(t, err, ErrBuildingProxyURL)
	})
}
