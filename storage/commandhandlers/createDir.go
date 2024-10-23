package commandhandlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"storage/config"
)

func CommandCreateDirHandler(w http.ResponseWriter, r *http.Request) {

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
	err = os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to create directories", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Directory created successfully"))
}
