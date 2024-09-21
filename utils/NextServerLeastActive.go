package utils

import "github.com/SicParv1sMagna/LoadBalancer/models"

func NextServerLeastActive(servers []*models.Server) *models.Server {
	leastActiveConnections := -1

	leastActiveServer := servers[0]

	for _, server := range servers {
		server.Mutex.Lock()

		if (server.ActiveConnections < leastActiveConnections || leastActiveConnections == -1) && server.Healthy {
			leastActiveConnections = server.ActiveConnections
			leastActiveServer = server
		}

		server.Mutex.Unlock()
	}

	return leastActiveServer
}
