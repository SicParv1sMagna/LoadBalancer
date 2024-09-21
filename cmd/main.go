package main

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/SicParv1sMagna/LoadBalancer/cmd/config"
	"github.com/SicParv1sMagna/LoadBalancer/models"
	"github.com/SicParv1sMagna/LoadBalancer/utils"
)

func main() {
	config := config.MustLoad()

	healthCheckInterval, err := time.ParseDuration(config.HealthCheckInterval)
	if err != nil {
		log.Fatalf("Invalid health check interval: %s", err.Error())
	}

	var servers []*models.Server

	for _, serverUrl := range config.Servers {
		u, _ := url.Parse(serverUrl)
		servers = append(servers, &models.Server{URL: u})
	}

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server := utils.NextServerLeastActive(servers)
		server.Mutex.Lock()
		server.ActiveConnections++
		server.Mutex.Unlock()
		server.Proxy().ServeHTTP(w, r)
		server.Mutex.Lock()
		server.ActiveConnections--
		server.Mutex.Unlock()
	})

	log.Println("Starting server on port", config.ListenPort)
	err = http.ListenAndServe(config.ListenPort, nil)
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}
