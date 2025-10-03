package api

import (
	"encoding/json"
	"net/http"
)

// Public API: accepts only N (number of users)
type StartRequest struct {
	N int `json:"n"`
}

type StartResponse struct {
	Started  bool   `json:"started"`
	WinnerID string `json:"winnerId,omitempty"`
	Round    int    `json:"round,omitempty"`
}

// Internal answer submission (used by simulator)
type AnswerRequest struct {
	UserID  string `json:"userId"`
	Correct bool   `json:"correct"`
}

type AnswerResponse struct {
	Accepted bool   `json:"accepted"`
	WinnerID string `json:"winnerId,omitempty"`
	Round    int    `json:"round,omitempty"`
}

func DecodeStart(r *http.Request) (StartRequest, error) {
	var req StartRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func DecodeAnswer(r *http.Request) (AnswerRequest, error) {
	var req AnswerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func EncodeJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
