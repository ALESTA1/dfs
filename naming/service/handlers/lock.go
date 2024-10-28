package handlers

import (
	"encoding/json"
	"io"
	"naming/config"
	"naming/directree"
	"naming/replication"
	"net/http"
	"strings"
)

func arrayToMap(arr []string) map[string]bool {
	result := make(map[string]bool)
	for _, str := range arr {
		result[str] = true
	}
	return result
}

func Lock(w http.ResponseWriter, r *http.Request) {

	
	type Body struct {
		Path      string `json:"path"`
		Exclusive bool   `json:"exclusive"`
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var body Body
	err = json.Unmarshal(b, &body)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	path := strings.Split(body.Path, "/")
	f := directree.IsValidPath(config.Root, 0, path)

	if f {

		directree.Lock(config.Root, 0, path, body.Exclusive)
		if !body.Exclusive {
			replicate, hosts, node := directree.CheckReplication(config.Root, 0, path)
			if replicate {
				hostsMap := arrayToMap(hosts)
				host := "none"

				for key := range config.StorageCommandPorts {

					_, exists := hostsMap[key]
					if !exists {
						host = key
						break
					}
				}

				if host != "none" {

					directree.Lock(config.Root, 0, path, false)
					go replication.Replicate(host, hosts[0], node, path, body.Path)
				}

			}
		} else {

			node := directree.CheckDereplication(config.Root, 0, path)

			if node != nil {

				replication.Dereplicate(node, body.Path)
			}
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}
