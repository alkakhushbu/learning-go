package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"os"
)

// The main function is the entry point of this application.
func main() {
	// Register the service with Consul for service discovery.
	registerServiceConsul()

	// Create a new Gin router object. `gin.Default()` sets up a router with default middleware.
	r := gin.Default()

	r.GET("/user/ping", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "pong from micro-1",
		})
	})

	panic(r.Run(":8089"))
}

// registerServiceConsul registers the HTTP service with Consul for service discovery.
func registerServiceConsul() {
	// Create a default configuration for connecting to Consul.
	config := api.DefaultConfig()

	// Set the address of the Consul server. Change this to point to your actual Consul service.
	config.Address = "http://consul.app:8500"

	// Create a new client to interact with Consul.
	consul, err := api.NewClient(config)
	if err != nil {
		// If an error occurs while creating the client, stop the application by panicking.
		panic(err)
	}

	// Fetch the hostname of the machine/service from the environment variable `HOSTNAME`.
	address := os.Getenv("HOSTNAME")
	if address == "" {
		// If `HOSTNAME` is not set, stop the application with an error message.
		panic("HOSTNAME not set")
	}

	// Create a new Consul service registration object.
	registration := &api.AgentServiceRegistration{}

	// Assign a name for the service. This is how other applications will reference this service in Consul.
	registration.Name = "micro-1"

	// Assign a unique ID for the service, combining the service name and the hostname to make it unique.
	registration.ID = "micro-1-" + address

	// Assign the hostname (service's address) and port on which this service is running.
	registration.Address = address
	registration.Port = 8089

	// Register the service with Consul. If this process fails, stop the application.
	err = consul.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	err = consul.Agent().ServiceDeregister(fmt.Sprintf("%s-%s", "micro-1-", address))
	if err != nil {
		panic(err)
	}
}
