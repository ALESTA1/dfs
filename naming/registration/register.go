package registration

import (
	"encoding/json"
	"fmt"
	"io"
	"naming/config"
	"naming/directree"
	"net/http"
	"strings"
)

type RegisterBody struct {
	Storage_ip   string   `json:"storage_ip"`
	Client_port  int      `json:"client_port"`
	Command_port int      `json:"command_port"`
	Files        []string `json:"files"`
}

func validateRequestBody(r *http.Request) (*RegisterBody, error) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %v", err)
	}

	var rawBody map[string]interface{}
	err = json.Unmarshal(body, &rawBody)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON format: %v", err)
	}

	expectedFields := map[string]struct{}{
		"storage_ip":   {},
		"client_port":  {},
		"command_port": {},
		"files":        {},
	}

	for key := range rawBody {
		if _, ok := expectedFields[key]; !ok {
			return nil, fmt.Errorf("unexpected field: %s", key)
		}
	}

	var regbody RegisterBody
	err = json.Unmarshal(body, &regbody)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON format: %v", err)
	}

	if regbody.Storage_ip == "" {
		return nil, fmt.Errorf("storage_ip is required")
	}
	if regbody.Client_port <= 0 {
		return nil, fmt.Errorf("client_port must be a positive integer")
	}
	if regbody.Command_port <= 0 {
		return nil, fmt.Errorf("command_port must be a positive integer")
	}

	return &regbody, nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	b, err := validateRequestBody(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	_, exists := config.StorageClientPorts[b.Storage_ip]
	if exists {
		http.Error(w, "Invalid request already registered storage server", http.StatusConflict)
		return
	}

	config.StorageClientPorts[b.Storage_ip] = b.Client_port
	config.StorageCommandPorts[b.Storage_ip] = b.Command_port
	filesToDelete := []string{}

	fmt.Println(b.Storage_ip)
	for _, file := range b.Files {

		temp := strings.Split(file, "/")
		fmt.Println(temp)
		isInserted := directree.Insert(config.Root, 0, temp, b.Storage_ip)
		fmt.Println(isInserted)
		if !isInserted {
			filesToDelete = append(filesToDelete, file)
		}
	}
	fmt.Println(filesToDelete)
	response := struct {
		Files []string `json:"files"`
	}{
		Files: filesToDelete,
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON response: %s", err.Error()), http.StatusInternalServerError)
	}

}
