package commandhandlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"storage/config"
)

func CommandCreateHandler(w http.ResponseWriter, r *http.Request) {

	type CreateRequest struct {
		Path string `json:"path"`
	}

	var req CreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(config.Directory, req.Path)
	fullPath = ".\\" + fullPath
	dir := filepath.Dir(fullPath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to create directories", http.StatusInternalServerError)
		return
	}

	// Create the file
	file, err := os.Create(fullPath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File created successfully"))
}
