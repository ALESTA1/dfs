package commandhandlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"storage/config"
)

func CommandDeleteHandler(w http.ResponseWriter, r *http.Request) {
	type DeleteRequest struct {
		Path string `json:"path"`
	}

	var req DeleteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(config.Directory, req.Path)
	fullPath = ".\\" + fullPath
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	err = os.Remove(fullPath)
	if err != nil {
		http.Error(w, "Error deleting file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"File deleted successfully"}`))
}
