package handlers

import (
	"encoding/json"
	"io"
	"naming/config"
	"naming/directree"
	"net/http"
	"strings"
)

func IsDir(w http.ResponseWriter, r *http.Request) {
	config.GlobalMutex.Lock()
	defer config.GlobalMutex.Unlock()
	type Body struct {
		Path string `json:"path"`
	}

	b, _ := io.ReadAll(r.Body)

	var body Body
	json.Unmarshal(b, &body)

	path := strings.Split(body.Path, "/")

	f := directree.IsDir(config.Root, 0, path)

	response := map[string]bool{
		"success": f == 0,
	}
	if f == 0 {
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else if f == 1 {
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}
