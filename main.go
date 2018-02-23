package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var kongHost = flag.String("kong", "", "kong host")
var kongAdminPort = flag.String("kong-admin-port", "8001", "kong admin port")

var healthCheckInterval = flag.String("health-check-interval", "2000", "healt check interval in ms")
var healthCheckPath = flag.String("health-check-path", "/ping", "path to check for active health check")

func main() {
	flag.Parse()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if *kongHost == "" {
		log.Fatalf("`kong` flag did not provide kong host")
	}

	kongClient := newKongClient(*kongHost, *kongAdminPort)
	kongHealthCheckConfig := &kongHealthCheckConfig{
		healthCheckPath:     *healthCheckPath,
		healthCheckInterval: *healthCheckInterval,
	}

	//TODO: make size of the targetChan configurable
	healthCheck := &kongHealthCheck{
		client:                kongClient,
		kongHealthCheckConfig: kongHealthCheckConfig,
		targetChan:            make(chan target, 100),
	}

	healthCheck.start()
	defer healthCheck.stop()

	log.Printf("started kong-healthcheck for kong host: %s with interval: %s ms", *kongHost, *healthCheckInterval)
	sig := <-sigChan
	log.Printf("stopping kong-healthcheck, received os signal: %v", sig)
}
