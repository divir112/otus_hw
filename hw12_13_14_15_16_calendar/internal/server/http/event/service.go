package event

import (
	"context"
	"net/http"

	"github.com/divir112/otus_hw/internal/model"
)

type Application interface {
	CreateEvent(ctx context.Context, event model.Event) error
	GetEvents(ctx context.Context) ([]model.Event, error)
	UpdateEvents(ctx context.Context, id int, event model.Event) error
	DeleteEvent(ctx context.Context, id int) error
	GetEventsByDate(ctx context.Context, days int, date string) ([]model.Event, error)
}

type Logger interface {
	Info(string, ...any)
	Error(string, ...any)
}

type EventServiceHTTP struct {
	app    Application
	logger Logger
	mux    *http.ServeMux
}

func NewEventServiceHTTP(app Application, logger Logger, mux *http.ServeMux) *EventServiceHTTP {
	return &EventServiceHTTP{
		app:    app,
		logger: logger,
		mux:    mux,
	}
}

func (s *EventServiceHTTP) Register() {
	s.mux.HandleFunc("GET /events", s.GetEvents)
	s.mux.HandleFunc("POST /events", s.CreateEvent)
	s.mux.HandleFunc("PUT /events/:id", s.UpdateEvent)
	s.mux.HandleFunc("DELETE /events/:id", s.DeleteEvent)
}
