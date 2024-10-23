package commandhandlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"storage/config"
	"strconv"
)

func CommandCopyHandler(w http.ResponseWriter, r *http.Request) {

	type CopyRequest struct {
		Path        string `json:"path"`
		Server_ip   string `json:"server_ip"`
		Server_port int    `json:"server_port"`
	}

	var req CopyRequest
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

	file, err := os.Create(fullPath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusNotFound)
		return
	}
	defer file.Close()

	size, err := getFileSize(req.Server_ip, req.Server_port, req.Path)
	if err != nil {
		http.Error(w, "Failed to get file size", http.StatusInternalServerError)
		return
	}

	fileData, err := readFileFromServer(req.Server_ip, req.Server_port, req.Path, size)
	if err != nil {
		http.Error(w, "Failed to read file data", http.StatusInternalServerError)
		return
	}

	_, err = file.Write(fileData)
	if err != nil {
		http.Error(w, "Failed to write to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File copied successfully"))
}

func getFileSize(serverIP string, serverPort int, path string) (int64, error) {
	url := "http://" + serverIP + ":" + strconv.Itoa(serverPort) + "/storage_size"

	reqBody, _ := json.Marshal(map[string]string{
		"path": path,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, err
	}

	var sizeResponse struct {
		Size int64 `json:"size"`
	}
	err = json.NewDecoder(resp.Body).Decode(&sizeResponse)
	if err != nil {
		return 0, err
	}

	return sizeResponse.Size, nil
}

func readFileFromServer(serverIP string, serverPort int, path string, size int64) ([]byte, error) {
	url := "http://" + serverIP + ":" + strconv.Itoa(serverPort) + "/storage_read"

	reqBody, _ := json.Marshal(map[string]interface{}{
		"path":   path,
		"offset": 0,
		"length": size,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	fileData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}
