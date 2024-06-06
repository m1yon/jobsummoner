package linkedin

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpProxyBuilder(t *testing.T) {
	fakeProxyConfig := ProxyConfig{
		Hostname: "crazyproxies.com",
		Port:     "44443",
		Username: "michael",
		Password: "hunter12",
	}

	url, err := BuildHttpProxyURL(fakeProxyConfig)
	assert.NoError(t, err)
	assert.Equal(t, "http://michael:hunter12@crazyproxies.com:44443", url.String())
}

func TestHttpProxyClient(t *testing.T) {
	t.Run("initializes proxy correctly", func(t *testing.T) {
		proxyCalled := false
		proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxyCalled = true
			w.WriteHeader(http.StatusOK)
		}))
		defer proxyServer.Close()

		proxyServerURL, err := url.Parse(proxyServer.URL)
		assert.NoError(t, err)

		client, err := NewHttpProxyClient(proxyServerURL)
		assert.NoError(t, err)

		req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
		_, err = client.Transport.RoundTrip(req)
		assert.NoError(t, err)
		assert.Equal(t, true, proxyCalled, "proxy was not called")
	})
}
