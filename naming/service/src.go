package service

import (
	"fmt"
	"naming/config"
	"naming/service/handlers"
	"net/http"
)

func ClientService() {

	http.HandleFunc("/is_valid_path", handlers.ValidPath)
	http.HandleFunc("/getstorage", handlers.GetStorage)
	http.HandleFunc("/delete", handlers.Delete)
	http.HandleFunc("/create_directory", handlers.CreateDir)
	http.HandleFunc("/create_file", handlers.CreateFile)
	http.HandleFunc("/is_directory", handlers.IsDir)
	http.HandleFunc("/unlock", handlers.Unlock)
	http.HandleFunc("/lock", handlers.Lock)
	http.HandleFunc("/list", handlers.List)

	fmt.Println("Starting client server on " + config.SERVICE_PORT)
	err := http.ListenAndServe(":"+config.SERVICE_PORT, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
