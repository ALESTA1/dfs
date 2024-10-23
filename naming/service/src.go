package service

import (
	"fmt"
	"naming/config"
	"naming/service/handlers"
	"net/http"
)

type Request struct {
	w http.ResponseWriter
	r *http.Request
}

var requestQueue = make(chan Request)

func processRequests() {
	for req := range requestQueue {
		switch req.r.URL.Path {
		case "/is_valid_path":
			handlers.ValidPath(req.w, req.r)
		case "/getstorage":
			handlers.GetStorage(req.w, req.r)
		case "/delete":
			handlers.Delete(req.w, req.r)
		case "/create_directory":
			handlers.CreateDir(req.w, req.r)
		case "/create_file":
			handlers.CreateFile(req.w, req.r)
		case "/is_directory":
			handlers.IsDir(req.w, req.r)
		case "/unlock":
			handlers.Unlock(req.w, req.r)
		case "/lock":
			handlers.Lock(req.w, req.r)
		case "/list":
			handlers.List(req.w, req.r)
		default:
			http.NotFound(req.w, req.r)
		}
	}
}

func ClientService() {

	go processRequests()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestQueue <- Request{w: w, r: r}
	})

	fmt.Println("Starting client server on " + config.SERVICE_PORT)
	err := http.ListenAndServe(":"+config.SERVICE_PORT, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
