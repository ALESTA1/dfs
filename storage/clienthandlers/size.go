package clienthandlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"storage/config"
)

func ClientSizeHandler(w http.ResponseWriter, r *http.Request) {
	type SizeRequest struct {
		Path string `json:"path"`
	}

	type SizeResponse struct {
		Size int64 `json:"size"`
	}

	var req SizeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(config.Directory, req.Path)
	fullPath = ".\\" + fullPath

	fileInfo, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error checking file", http.StatusInternalServerError)
		return
	}

	if fileInfo.IsDir() {
		http.Error(w, "Path is a directory, not a file", http.StatusBadRequest)
		return
	}

	fileSize := fileInfo.Size()

	response := SizeResponse{
		Size: fileSize,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error generating JSON response", http.StatusInternalServerError)
		return
	}
}
