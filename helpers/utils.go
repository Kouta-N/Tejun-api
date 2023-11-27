package helpers

import (
	"encoding/json"
	"net/http"
)

func WriteErr(w http.ResponseWriter, statusCode int, err error) {
	body, _ := json.Marshal(err)
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}