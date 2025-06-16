package res

import (
	"encoding/json"
	"net/http"
)

const StatusError = "Error"

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func Json(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
