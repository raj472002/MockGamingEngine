package game

import (
	"context"
	"encoding/json"
	"net/http"

	kittransport "github.com/go-kit/kit/transport/http"
)

func MakeHTTPHandler(eps Endpoints) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/submit", kittransport.NewServer(
		eps.Submit,
		decodeSubmit,
		encodeJSON,
	))
	return mux
}

func decodeSubmit(_ context.Context, r *http.Request) (interface{}, error) {
	var req submitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeJSON(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
