package handlers

import (
	"net/http"
	"time"
)

// Info - a dummy route handler that isn't protected by authentication
func Info(w http.ResponseWriter, r *http.Request) {
	message := map[string]string{
		"status":      "awesome",
		"currentTime": time.Now().Format("2016-01-02 15:04:05"),
	}
	SendResponse(w, message, 200)
}
