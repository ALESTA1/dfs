package clienthandlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"storage/config"
)

type WriteRequest struct {
	Path   string `json:"path"`
	Offset int64  `json:"offset"`
	Data   string `json:"data"`
}

type BooleanReturn struct {
	Success bool `json:"success"`
}

func ClientWriteHandler(w http.ResponseWriter, r *http.Request) {
	var writeReq WriteRequest
	err := json.NewDecoder(r.Body).Decode(&writeReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(config.Directory, writeReq.Path)
	fullPath = ".\\" + fullPath
	if writeReq.Offset < 0 {
		http.Error(w, `{"exception_type": "IndexOutOfBoundsException", "exception_info": "Offset cannot be negative."}`, http.StatusBadRequest)
		return
	}

	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, `{"exception_type": "FileNotFoundException", "exception_info": "File not found."}`, http.StatusNotFound)
		return
	}
	defer file.Close()

	dataBytes, err := base64.StdEncoding.DecodeString(writeReq.Data)
	if err != nil {
		http.Error(w, "Invalid base64 data", http.StatusBadRequest)
		return
	}

	_, err = file.Seek(writeReq.Offset, 0)
	if err != nil {
		http.Error(w, `{"exception_type": "IOException", "exception_info": "Error seeking in file."}`, http.StatusInternalServerError)
		return
	}

	_, err = file.Write(dataBytes)
	if err != nil {
		http.Error(w, `{"exception_type": "IOException", "exception_info": "Error writing to file."}`, http.StatusInternalServerError)
		return
	}

	response := BooleanReturn{
		Success: true,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
