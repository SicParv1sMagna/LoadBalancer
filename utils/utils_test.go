package utils_test

import (
	"net/url"
	"sync"
	"testing"

	"github.com/SicParv1sMagna/LoadBalancer/models"
	"github.com/SicParv1sMagna/LoadBalancer/utils"
)

func TestNextServerLeastActive(t *testing.T) {
	url1, _ := url.Parse("http://localhost:8081")
	url2, _ := url.Parse("http://localhost:8082")

	server1 := &models.Server{
		URL:               url1,
		ActiveConnections: 5,
		Healthy:           true,
		Mutex:             sync.Mutex{},
	}
	server2 := &models.Server{
		URL:               url2,
		ActiveConnections: 3,
		Healthy:           true,
		Mutex:             sync.Mutex{},
	}

	servers := []*models.Server{server1, server2}

	leastActiveServer := utils.NextServerLeastActive(servers)
	if leastActiveServer != server2 {
		t.Errorf("Expected server2 to be least active, got %v", leastActiveServer)
	}

	server2.Healthy = false

	leastActiveServer = utils.NextServerLeastActive(servers)
	if leastActiveServer != server1 {
		t.Errorf("Expected server1 to be selected, got %v", leastActiveServer)
	}
}
