# Load Balancer

### Requirements:

1. The load balancer should handle multiple incoming connections concurrently.
2. The load balancer should distribute the load to the least loaded server using the Least Connection algorithm.
3. The load balancer needs to implement healthchecks and only route requests to healthy backend servers.

### Config.yaml

First of all, you need to define ```./config/config.yaml```.

Example:
```yaml
env: "local"
listenPort: ":8080"
healthCheckInterval: "5s"
servers:
  - "http://localhost:8081"
  - "http://localhost:8082"
  - "http://localhost:8083"
```

Explanation:

```env``` - environment of balancer.
```listenPort``` - launch port.
```healthCheckInterval``` - interval between server healthcheck.
```servers``` - list of servers for distributing load.

### How to launch?

After creating of config.yaml you need to use this command for launch
```zsh
go run cmd/main.go --config=./config/config.yaml
```