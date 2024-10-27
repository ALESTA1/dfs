package main

import (
	"fmt"
	"naming/config"
	"naming/registration"
	"naming/service"
	"os"
	"time"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Please provide at least 2 ports.")
		return
	}
	config.IP = config.ResolveHostIp()
	println("Naming server running on IP .. " + config.IP)
	config.SERVICE_PORT = os.Args[1]
	config.REGISTRATION_PORT = os.Args[2]

	go service.ClientService()
	go registration.RegistrationService()

	for {
		time.Sleep(1000000000 * time.Second)
	}
}
