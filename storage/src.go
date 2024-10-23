package main

import (
	"fmt"
	"net/http"
	"os"
	"storage/clienthandlers"
	"storage/commandhandlers"
	"storage/config"
	"time"
)

func startCommandServer() {

	Register()
	http.HandleFunc("/storage_create", commandhandlers.CommandCreateFileHandler)
	http.HandleFunc("/storage_delete",commandhandlers.CommandDeleteHandler)
	http.HandleFunc("/storage_copy", commandhandlers.CommandCopyHandler)
	http.HandleFunc("/storage_create_dir", commandhandlers.CommandCreateDirHandler)

	fmt.Println("Starting storage command server at port " + config.COMMAND_PORT)
	println(config.IP+":"+config.COMMAND_PORT)
	err := http.ListenAndServe(config.IP+":"+config.COMMAND_PORT, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}


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
