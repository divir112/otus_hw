package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/divir112/otus_hw/internal/config"
	"github.com/divir112/otus_hw/internal/model"
	"github.com/divir112/otus_hw/internal/server/http/event"
)

type Server struct {
	Mux         *http.ServeMux
	server      *http.Server
	Logger      Logger
	Application Application
	Config      *config.Config
}

type Logger interface {
	Info(string, ...any)
	Error(string, ...any)
}

type Application interface {
	CreateEvent(ctx context.Context, event model.Event) error
	GetEvents(ctx context.Context) ([]model.Event, error)
	UpdateEvents(ctx context.Context, id int, event model.Event) error
	DeleteEvent(ctx context.Context, id int) error
	GetEventsByDate(ctx context.Context, days int, date string) ([]model.Event, error)
}

type Storage interface {
	Add(context.Context, model.Event) (int, error)
	Update(context.Context, int, model.Event) error
	Delete(context.Context, int) error
	List(context.Context) ([]model.Event, error)
}

func NewServer(logger Logger, app Application, config *config.Config) *Server {
	mux := http.NewServeMux()
	mw := NewMiddlware(logger)
	eventControllers := event.NewEventServiceHTTP(app, logger, mux)
	eventControllers.Register()
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mw.loggingMiddleware(mux),
		ReadHeaderTimeout: 15 * time.Second,
	}
	return &Server{mux, server, logger, app, config}
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
