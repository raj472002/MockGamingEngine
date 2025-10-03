package sse

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func Stream(ctx context.Context, w http.ResponseWriter, events <-chan interface{}) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-events:
			if !ok {
				return
			}
			w.Write([]byte("data: "))
			_ = enc.Encode(evt)
			w.Write([]byte("\n"))
			flusher.Flush()
		case <-ticker.C:
			w.Write([]byte(": keep-alive\n\n"))
			flusher.Flush()
		}
	}
}
