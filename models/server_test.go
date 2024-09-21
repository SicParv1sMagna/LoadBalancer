package models_test

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/SicParv1sMagna/LoadBalancer/models"
	"github.com/stretchr/testify/assert"
)

func TestServer_Proxy(t *testing.T) {
	serverUrl, err := url.Parse("http://localhost:8081")
	assert.NoError(t, err)

	server := &models.Server{
		URL: serverUrl,
	}

	proxy := server.Proxy()
	assert.NotNil(t, proxy, "Expected Proxy() to return a non-nil ReverseProxy")

	req := httptest.NewRequest("GET", "/", nil)

	proxyDirector := proxy.Director
	assert.NotNil(t, proxyDirector, "Expected ReverseProxy Director to be set")

	proxyDirector(req)

	assert.Equal(t, "localhost:8081", req.URL.Host, "Expected request to be directed to server URL")
	assert.Equal(t, "http", req.URL.Scheme, "Expected request scheme to match server URL scheme")
}
