package main

import (
	"net/http"
)

// StatusRetrieve returns a response for /payloads/statuses
type StatusRetrieve struct {
	RequestID string `json:"request_id"`
	Status string `json:"status"`
	ID string `json:"id"`
	Service string `json:"service"`
	Source string `json:"source"`
	StatusMsg string `json:"status_msg"`
	Date string `json:"date"`
	CreatedAt string `json:"created_at"`
}

func statuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Statuses"))
}
