package main

import (
	"fmt"
	"net/http"
	"os"
	"storage/clienthandlers"
	"storage/config"
	"time"
)

func startCommandServer() {

	Register()

}

func startClientServer() {

	http.HandleFunc("/storage_size", clienthandlers.ClientSizeHandler)
	http.HandleFunc("/storage_read", clienthandlers.ClientReadHandler)
	http.HandleFunc("/storge_write", clienthandlers.ClientWriteHandler)

	fmt.Println("Starting storage client server at port " + config.CLIENT_PORT)
	err := http.ListenAndServe(":"+config.CLIENT_PORT, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func main() {

	if len(os.Args) < 4 {
		fmt.Println("Please provide at least 2 ports.")
		return
	}

	config.CLIENT_PORT = os.Args[1]
	config.COMMAND_PORT = os.Args[2]
	config.REGISTRATION_PORT = os.Args[3]
	config.Directory = os.Args[4]
	os.MkdirAll(config.Directory, os.ModePerm)

	go startCommandServer()
	go startClientServer()

	for {
		time.Sleep(1000000000 * time.Second)
	}

}
