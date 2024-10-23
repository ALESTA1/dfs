package main

import (
	"fmt"
	"naming/config"
	"naming/registration"
	"naming/service"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Please provide at least 2 ports.")
		return
	}
	config.SERVICE_PORT = os.Args[1]
	config.REGISTRATION_PORT = os.Args[2]

	go service.ClientService()
	go registration.RegistrationService()

}
