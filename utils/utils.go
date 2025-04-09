package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

func UnmarshalJson(data []byte, res map[string]interface{}) error {
	if err := json.Unmarshal(data, &res); err != nil {
		return fmt.Errorf("failed to unmarshal json: %v", err)
	}
	return nil
}


func ApiKeyMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-API-Key")
        if apiKey == "" {
            ErrorResponse(w, http.StatusUnauthorized, 
                        fmt.Errorf("API key is missing"))
            return
        }

        validApiKey := os.Getenv("API_KEY")
        if apiKey != validApiKey {
            ErrorResponse(w, http.StatusUnauthorized, 
                        fmt.Errorf("invalid API key"))
            return
        }

        next.ServeHTTP(w, r)
    }
}