package clienthandlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"storage/config"
)

type ReadRequest struct {
	Path   string `json:"path"`
	Offset int64  `json:"offset"`
	Length int    `json:"length"`
}

type DataReturn struct {
	Data string `json:"data"`
}

func ClientReadHandler(w http.ResponseWriter, r *http.Request) {

	var readReq ReadRequest
	err := json.NewDecoder(r.Body).Decode(&readReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(config.Directory, readReq.Path)
	fullPath = ".\\" + fullPath

	file, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	_, err = file.Seek(readReq.Offset, 0)
	if err != nil {
		http.Error(w, "Error seeking in file", http.StatusInternalServerError)
		return
	}

	buffer := make([]byte, readReq.Length)
	bytesRead, err := file.Read(buffer)
	if err != nil {
		http.Error(w, "Error reading from file", http.StatusInternalServerError)
		return
	}

	encodedData := base64.StdEncoding.EncodeToString(buffer[:bytesRead])

	response := DataReturn{
		Data: encodedData,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
