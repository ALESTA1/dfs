package handlers

import (
	"encoding/json"
	"io"
	"naming/config"
	"naming/directree"
	"net/http"
	"strings"
)

func ValidPath(w http.ResponseWriter, r *http.Request) {

	type Body struct {
		Path string `json:"path"`
	}

	b, _ := io.ReadAll(r.Body)

	var body Body
	json.Unmarshal(b, &body)

	path := strings.Split(body.Path, "/")

	f := directree.IsValidPath(config.Root, 0, path)

	response := map[string]bool{
		"success": f,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}
