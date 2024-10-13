package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

var root *Node
var storageClientPorts = make(map[string]int)
var storageCommandPorts = make(map[string]int)
var globalMutex sync.Mutex
var REGISTRATION_PORT string
var SERVICE_PORT string

type RegisterBody struct {
	Storage_ip   string   `json:"storage_ip"`
	Client_port  int      `json:"client_port"`
	Command_port int      `json:"command_port"`
	Files        []string `json:"files"`
}

func validateRequestBody(r *http.Request) (*RegisterBody, error) {
	globalMutex.Lock()
	defer globalMutex.Unlock()

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

func registerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	b, err := validateRequestBody(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	_, exists := storageClientPorts[b.Storage_ip]
	if exists {
		http.Error(w, "Invalid request already registered storage server", http.StatusConflict)
		return
	}

	storageClientPorts[b.Storage_ip] = b.Client_port
	storageCommandPorts[b.Storage_ip] = b.Command_port
	filesToDelete := []string{}

	fmt.Println(b.Storage_ip)
	for _, file := range b.Files {

		temp := strings.Split(file, "/")
		fmt.Println(temp)
		isInserted := Insert(root, 0, temp, b.Storage_ip)
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
func main() {

	if len(os.Args) < 3 {
		fmt.Println("Please provide at least 2 ports.")
		return
	}
	SERVICE_PORT = os.Args[1]
	REGISTRATION_PORT = os.Args[2]

	http.HandleFunc("/register", registerHandler)
	root = NewNode("root")
	fmt.Println("Starting registration server on" + REGISTRATION_PORT + "... for Storage registration")
	err := http.ListenAndServe(":"+REGISTRATION_PORT, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}

type Node struct {
	Name     string
	Children map[string]*Node
	Hosts    []string
}

func NewNode(name string) *Node {
	return &Node{
		Name:     name,
		Children: make(map[string]*Node),
		Hosts:    []string{},
	}
}

func Insert(node *Node, i int, path []string, host string) bool {

	if i == len(path) {
		return true
	}
	currentNode := path[i]

	nextNode, exists := node.Children[currentNode]
	f := true
	if !exists {
		temp := NewNode(currentNode)
		node.Children[currentNode] = temp
		f = f && Insert(temp, i+1, path, host)
	} else {

		if i == len(path)-1 {
			return false
		}
		f = f && Insert(nextNode, i+1, path, host)
	}
	if f {
		node.Hosts = append(node.Hosts, host)
	}
	return f
}
