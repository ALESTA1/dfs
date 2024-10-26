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

func CreateFile(w http.ResponseWriter, r *http.Request) {
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

	if !f {
		host := config.GetRandomKey(config.StorageCommandPorts)
		directree.Insert(config.Root, 0, path, host)

		reqBody := struct {
			Path string `json:"path"`
		}{Path: body.Path}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			http.Error(w, "Error creating JSON body", http.StatusInternalServerError)
			return
		}
		println("http://"+host+":"+strconv.Itoa(config.StorageCommandPorts[host])+"/storage_create")
		resp, err := http.Post("http://"+host+":"+strconv.Itoa(config.StorageCommandPorts[host])+"/storage_create", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			http.Error(w, "Error sending POST request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Error from command server", resp.StatusCode)
			return
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
