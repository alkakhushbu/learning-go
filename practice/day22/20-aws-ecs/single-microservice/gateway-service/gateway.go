package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

var routeMap = map[string]string{
	"/hello/service": "hello-service",
}

func main() {
	router := gin.Default()

	router.GET("/hello/service", HelloHandler)
	router.Run(":80")
}

func HelloHandler(c *gin.Context) {
	path := c.Request.URL.Path
	svcName, ok := routeMap[path]
	if !ok {
		log.Println("URL path not found:", path)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "URL path not found:"})
		return
	}

	proxyToService(svcName, c)
}

func proxyToService(serviceName string, c *gin.Context) {
	// Create a default configuration for Consul.
	config := api.DefaultConfig()

	// Setting the address where Consul is running. Change this to point to your actual Consul server.
	config.Address = "http://consul.app:8500"

	// Create a new client to interact with the Consul API.
	consul, err := api.NewClient(config)
	if err != nil {
		// If an error occurs while creating the Consul client, return a 500 error (Internal Server Error).
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	// Query Consul for the service with the given name.
	services, _, err := consul.Health().Service(serviceName, "", true, nil)
	if err != nil {
		// If an error occurs while querying Consul, return a 500 error.
		c.JSON(http.StatusServiceUnavailable, gin.H{"Error": "Service Unavailable"})
		return
	}
	// If no services are found, return a 501 error.
	if len(services) == 0 {
		c.JSON(http.StatusNotImplemented, gin.H{"Error": "Service not found"})
		return
	}

	// Pick the first available service instance (can be enhanced later for load balancing).
	service := services[0]
	log.Printf("Service :%s, port: %d\n", service.Service.Address, service.Service.Port)

	// Construct the URL to forward the request to the service.
	// `ServiceAddress` and `ServicePort` are the address and port of the service found in Consul.
	serviceAddress := fmt.Sprintf("http://%s:%d%s", service.Service.Address, service.Service.Port, c.Request.URL.Path)
	log.Printf("Service Address: %s\n", serviceAddress)

	// Make an HTTP GET request to the constructed service address.
	res, err := http.Get(serviceAddress)
	if err != nil {
		// If an error occurs while forwarding the request, return a 500 error.
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	// Read the response body from the service response.
	b, err := io.ReadAll(res.Body)
	if err != nil {
		// If an error occurs while reading the response body, return a 500 error.
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error reading response from " + serviceName})
		return
	}

	// Set the same response status code as that of the service response.
	// Forward the service response back to the requester (client).
	c.String(http.StatusOK, string(b))
}
