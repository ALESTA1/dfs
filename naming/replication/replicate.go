package replication

import (
	"bytes"
	"encoding/json"
	"naming/config"
	"naming/directree"
	"net/http"
	"strconv"
)

func Replicate(toHost string, fromHost string, node *directree.Node, path []string, originalPath string) {

	reqBody := struct {
		Path        string `json:"path"`
		Server_ip   string `json:"server_ip"`
		Server_port int    `json:"server_port"`
	}{Path: originalPath, Server_ip: fromHost, Server_port: config.StorageCommandPorts[fromHost]}

	jsonBody, _ := json.Marshal(reqBody)
	endpoint := "http://" + toHost + ":" + strconv.Itoa(config.StorageCommandPorts[toHost]) + "/storage_copy"
	http.Post(endpoint, "application/json", bytes.NewBuffer(jsonBody))
	//update hosts for current file
	node.HostsLock.Lock()
	node.Hosts = append(node.Hosts, toHost)
	node.HostsLock.Unlock()
	directree.Unlock(config.Root, 0, path, false)
}
