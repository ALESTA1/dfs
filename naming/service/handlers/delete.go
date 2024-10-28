package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"naming/config"
	"naming/directree"
	"net/http"
	"strconv"
	"strings"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	
	type Body struct {
		Path string `json:"path"`
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var body Body
	err = json.Unmarshal(b, &body)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	path := strings.Split(body.Path, "/")
	f := directree.IsValidPath(config.Root, 0, path)

	if f {
		directree.Lock(config.Root, 0, path, true)
		hosts := directree.Delete(config.Root, 0, path)
		directree.Unlock(config.Root, 0, path, true)
		for _, host := range hosts {

			reqBody := struct {
				Path string `json:"path"`
			}{Path: body.Path}

			jsonBody, err := json.Marshal(reqBody)
			if err != nil {
				http.Error(w, "Error creating JSON body", http.StatusInternalServerError)
				return
			}
			endpoint := "http://" + host + ":" + strconv.Itoa(config.StorageCommandPorts[host]) + "/storage_delete"
			println(endpoint)
			resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonBody))
			if err != nil {
				http.Error(w, "Error sending POST request", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				http.Error(w, "Error from storage server", resp.StatusCode)
				return
			}
		}
		response := map[string]bool{"success": true}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		response := map[string]bool{"success": false}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
	
}
