package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/divir112/otus_hw/internal/model"
)

type Server struct {
	Mux         *http.ServeMux
	server      *http.Server
	Logger      Logger
	Application Application
}

type Logger interface {
	Info(string, ...any)
}

type Application interface { // TODO
}

type Storage interface {
	Add(context.Context, model.Event) (int, error)
	Update(context.Context, int, model.Event) error
	Delete(context.Context, int) error
	List(context.Context) ([]model.Event, error)
}

func NewServer(logger Logger, app Application) *Server {
	mux := http.NewServeMux()
	mw := NewMiddlware(logger)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mw.loggingMiddleware(mux),
		ReadHeaderTimeout: 15 * time.Second,
	}
	return &Server{mux, server, logger, app}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("can't run server %w", err)
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("can't stop server %w", err)
	}
	return nil
}
