package game

import (
	"context"
	"log"
	"sync"
	"time"
)

type Submission struct {
	UserID  string
	Correct bool
	SentAt  time.Time
}

type Winner struct {
	UserID string `json:"userId"`
	Round  int    `json:"round"`
}

type Service interface {
	StartRound(ctx context.Context, round int, correctAnswer string)
	Submit(ctx context.Context, s Submission) (bool, *Winner)
	Subscribe() (<-chan Winner, func())
}

type service struct {
	mu              sync.Mutex
	currentRound    int
	correctAnswer   string
	winnerAnnounced bool

	subscribersMu sync.Mutex
	subscribers   map[chan Winner]struct{}
}

func NewService() Service {
	return &service{
		subscribers: make(map[chan Winner]struct{}),
	}
}

func (s *service) StartRound(ctx context.Context, round int, correctAnswer string) {
	s.mu.Lock()
	s.currentRound = round
	s.correctAnswer = correctAnswer
	s.winnerAnnounced = false
	s.mu.Unlock()
}

func (s *service) Submit(ctx context.Context, sub Submission) (bool, *Winner) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.winnerAnnounced {
		return false, nil
	}
	if sub.Correct {
		s.winnerAnnounced = true
		w := Winner{UserID: sub.UserID, Round: s.currentRound}
		log.Printf("winner declared: %s (round=%d)", w.UserID, w.Round)
		s.broadcast(w)
		return true, &w
	}
	return false, nil
}

func (s *service) Subscribe() (<-chan Winner, func()) {
	ch := make(chan Winner, 1)
	s.subscribersMu.Lock()
	s.subscribers[ch] = struct{}{}
	s.subscribersMu.Unlock()

	cancel := func() {
		s.subscribersMu.Lock()
		delete(s.subscribers, ch)
		s.subscribersMu.Unlock()
		close(ch)
	}
	return ch, cancel
}

func (s *service) broadcast(w Winner) {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()
	for ch := range s.subscribers {
		select {
		case ch <- w:
		default:
		}
	}
}
