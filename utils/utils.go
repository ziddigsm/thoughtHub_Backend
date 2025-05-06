package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
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

var rateLimiterMap sync.Map

func getRateLimit(apiKey string) *rate.Limiter {
	limiterMap, _ := rateLimiterMap.LoadOrStore(apiKey, rate.NewLimiter(15, 23)) //this is the burst(23) and rate limit(15)
	limiter := limiterMap.(*rate.Limiter)
	return limiter
}

func RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			ErrorResponse(w, http.StatusUnauthorized,
				fmt.Errorf("user not authorized"))
			return
		}
		limiter := getRateLimit(apiKey)

		if !limiter.Allow() {
			ErrorResponse(w, http.StatusTooManyRequests,
				fmt.Errorf("too many requests. please try again later"))
			return
		}
		next.ServeHTTP(w, r)
	}
}

func ApiKeyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			ErrorResponse(w, http.StatusUnauthorized,
				fmt.Errorf("user not authorized"))
			return
		}
		currentDate := time.Now().Weekday()
		apiKeyEnv := "API_KEY_" + strconv.Itoa(int(currentDate))
		validApiKey := os.Getenv(apiKeyEnv)
		if apiKey != validApiKey {
			ErrorResponse(w, http.StatusUnauthorized,
				fmt.Errorf("user unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	}
}
