package game

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type startRoundRequest struct {
	Round         int    `json:"round"`
	CorrectAnswer string `json:"correctAnswer"`
}

type startRoundResponse struct{}

type submitRequest struct {
	UserID  string `json:"userId"`
	Correct bool   `json:"correct"`
}

type submitResponse struct {
	Accepted bool   `json:"accepted"`
	WinnerId string `json:"winnerId,omitempty"`
	Round    int    `json:"round,omitempty"`
}

type Endpoints struct {
	Submit endpoint.Endpoint
}

func MakeEndpoints(svc Service) Endpoints {
	return Endpoints{
		Submit: makeSubmitEndpoint(svc),
	}
}

func makeSubmitEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(submitRequest)
		ok, w := svc.Submit(ctx, Submission{UserID: req.UserID, Correct: req.Correct})
		if w != nil {
			return submitResponse{Accepted: ok, WinnerId: w.UserID, Round: w.Round}, nil
		}
		return submitResponse{Accepted: ok}, nil
	}
}
