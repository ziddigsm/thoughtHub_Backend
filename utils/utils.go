package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseRequest(r *http.Request, reqBody interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	return json.NewDecoder(r.Body).Decode(reqBody)
}

func SuccessResponse(w http.ResponseWriter, status int, res any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(res)
}

func ErrorResponse(w http.ResponseWriter, status int, err error) {
	SuccessResponse(w, status, map[string]string{"message": err.Error()})
}