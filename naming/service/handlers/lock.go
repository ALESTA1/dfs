package handlers

import (
	"encoding/json"
	"io"
	"naming/config"
	"naming/directree"
	"net/http"
	"strings"
)

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
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
