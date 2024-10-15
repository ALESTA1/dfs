package service

import (
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

}
