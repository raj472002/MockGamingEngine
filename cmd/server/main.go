package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"gamingEngine/internal/api"
	"gamingEngine/internal/game"
)

func main() {
	logger := kitlog.NewLogfmtLogger(os.Stdout)
	_ = level.Info(logger).Log("msg", "starting server")

	svc := game.NewService()
	// initialize a default round; simulator and clients send /submit answers
	svc.StartRound(context.Background(), 1, "correct")
	eps := game.MakeEndpoints(svc)
	gameMux := game.MakeHTTPHandler(eps)

	apiServer := api.NewServer(svc)
	mux := http.NewServeMux()
	mux.HandleFunc("/start", apiServer.StartHandler)
	mux.Handle("/", gameMux)
	http.Handle("/", mux)

	srv := &http.Server{Addr: ":8080", ReadHeaderTimeout: 5 * time.Second}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			_ = level.Error(logger).Log("err", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
