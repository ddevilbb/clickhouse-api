package helper

import (
	"encoding/json"
	"net/http"
)

func Respond(writer http.ResponseWriter, data map[string]interface{}, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	json.NewEncoder(writer).Encode(data)
}
