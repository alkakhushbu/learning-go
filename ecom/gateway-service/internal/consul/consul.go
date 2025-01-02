package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"math/rand"
	"time"
)

func GetService(client *consulapi.Client, serviceName string) (string, int, error) {

	services, _, err := client.Health().Service(serviceName, "", true, nil)
	if err != nil || len(services) == 0 {
		return "", 0, err
	}
	var service *consulapi.ServiceEntry
	fmt.Println(services)
	if len(services) > 1 {
		fmt.Println("more than one service")
		fmt.Printf("%+v\n", services)
		source := rand.NewSource(time.Now().UnixNano())
		rng := rand.New(source)
		randomServiceIndex := rng.Intn(len(services)) // 3 , 0-2
		service = services[randomServiceIndex]
	} else {
		// Select the only available instance
		service = services[0]
	}

	return service.Service.Address, service.Service.Port, nil

}
