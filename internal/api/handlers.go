package api

import (
	"gamingEngine/internal/game"
	"gamingEngine/internal/sim"
	"net/http"
	"time"
)

type Server struct {
	Game game.Service
}

func NewServer(gameSvc game.Service) *Server {
	return &Server{Game: gameSvc}
}

// StartHandler handles POST /start with {"n": <int>}
func (s *Server) StartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	req, err := DecodeStart(r)
	if err != nil || req.N <= 0 {
		http.Error(w, "invalid n", http.StatusBadRequest)
		return
	}
	// Start round with a fixed correct answer
	s.Game.StartRound(r.Context(), 1, "correct")
	ch, cancel := s.Game.Subscribe()
	defer cancel()
	go sim.Run("http://localhost:8080", req.N, "correct")
	// wait for winner or timeout
	select {
	case wnr := <-ch:
		_ = EncodeJSON(w, http.StatusOK, StartResponse{Started: true, WinnerID: wnr.UserID, Round: wnr.Round})
	case <-time.After(5 * time.Second):
		_ = EncodeJSON(w, http.StatusOK, StartResponse{Started: true})
	}
}

// AnswerHandler handles POST /answer for simulator submissions
func (s *Server) AnswerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	req, err := DecodeAnswer(r)
	if err != nil || req.UserID == "" {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	accepted, winner := s.Game.Submit(r.Context(), game.Submission{UserID: req.UserID, Correct: req.Correct})
	resp := AnswerResponse{Accepted: accepted}
	if winner != nil {
		resp.WinnerID = winner.UserID
		resp.Round = winner.Round
	}
	_ = EncodeJSON(w, http.StatusOK, resp)
}
