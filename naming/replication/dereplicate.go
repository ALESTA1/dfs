package replication

import (
	"bytes"
	"encoding/json"
	"naming/config"
	"naming/directree"
	"net/http"
	"strconv"
)

func Dereplicate(node *directree.Node, path string) {

	hosts := node.Hosts

	for i := 1; i < len(hosts); i++ {

		host := hosts[i]
		reqBody := struct {
			Path string `json:"path"`
		}{Path: path}

		jsonBody, _ := json.Marshal(reqBody)
		endpoint := "http://" + host + ":" + strconv.Itoa(config.StorageCommandPorts[host]) + "/storage_delete"
		http.Post(endpoint, "application/json", bytes.NewBuffer(jsonBody))

	}

	node.Hosts = node.Hosts[:1]

}
