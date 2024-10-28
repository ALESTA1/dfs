package handlers

import (
	"encoding/json"
	"io"
	"naming/config"
	"naming/directree"
	"net/http"
	"strings"
)

func GetStorage(w http.ResponseWriter, r *http.Request) {
	
	type Body struct {
		Path string `json:"path"`
	}

	b, _ := io.ReadAll(r.Body)

	var body Body
	json.Unmarshal(b, &body)

	path := strings.Split(body.Path, "/")

	f := directree.IsDir(config.Root, 0, path)

	if f == 1 {

		host := directree.GetHost(config.Root, 0, path)
		type Response struct {
			ServerIP   string `json:"server_ip"`
			ServerPort int    `json:"server_port"`
		}
		responseData := Response{
			ServerIP:   host,
			ServerPort: config.StorageClientPorts[host],
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(responseData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	
}
