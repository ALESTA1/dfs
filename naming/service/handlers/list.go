package handlers

import (
	"encoding/json"
	"io"
	"naming/config"
	"naming/directree"
	"net/http"
	"strings"
)

func List(w http.ResponseWriter, r *http.Request) {

	
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

	path := strings.Split(strings.Trim(body.Path, "/"), "/")
	node := directree.FindNode(config.Root, 0, path)

	if node != nil {
		var files []string

		tempPath := path
		tempPath = append(tempPath, "/")

		directree.Lock(config.Root, 0, tempPath, false)
		directree.List(node, "/"+node.Name, &files)
		directree.Unlock(config.Root, 0, tempPath, false)

		response := map[string][]string{
			"files": files,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid path"))
	}


}
