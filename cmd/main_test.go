package main_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/SicParv1sMagna/LoadBalancer/models"
	"github.com/SicParv1sMagna/LoadBalancer/utils"
)

func TestHealthCheck(t *testing.T) {
	healthyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer healthyServer.Close()

	unhealthyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer unhealthyServer.Close()

	uHealthy, _ := url.Parse(healthyServer.URL)
	uUnhealthy, _ := url.Parse(unhealthyServer.URL)

	servers := []*models.Server{
		{URL: uHealthy, Healthy: true},
		{URL: uUnhealthy, Healthy: true},
	}

	healthCheckInterval := 50 * time.Millisecond

	for _, server := range servers {
		go func(s *models.Server) {
			for range time.Tick(healthCheckInterval) {
				res, err := http.Get(s.URL.String())
				if err != nil || res.StatusCode >= 500 {
					s.Healthy = false
				} else {
					s.Healthy = true
				}
			}
		}(server)
	}

	time.Sleep(200 * time.Millisecond)

	if !servers[0].Healthy {
		t.Errorf("Expected server 1 to be healthy")
	}

	if servers[1].Healthy {
		t.Errorf("Expected server 2 to be unhealthy")
	}
}

func TestRequestForwarding(t *testing.T) {
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer targetServer.Close()

	u, _ := url.Parse(targetServer.URL)

	server := &models.Server{
		URL:               u,
		ActiveConnections: 0,
		Healthy:           true,
		Mutex:             sync.Mutex{},
	}

	servers := []*models.Server{server}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	serverToUse := utils.NextServerLeastActive(servers)
	serverToUse.Proxy().ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %d", resp.StatusCode)
	}
}
