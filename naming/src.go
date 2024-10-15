package main

import (
	"fmt"
	"naming/config"
	"naming/directree"
	"naming/registration"
	"net/http"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Please provide at least 2 ports.")
		return
	}
	config.SERVICE_PORT = os.Args[1]
	config.REGISTRATION_PORT = os.Args[2]

	http.HandleFunc("/register", registration.RegisterHandler)
	config.Root = directree.NewNode("root")
	fmt.Println("Starting registration server on" + config.REGISTRATION_PORT + " ... for Storage registration")
	err := http.ListenAndServe(":"+config.REGISTRATION_PORT, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}
