package event

import (
	"context"

	"github.com/divir112/otus_hw/internal/model"
	"github.com/divir112/otus_hw/internal/pb/event_service"
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

type EventService struct {
	app    Application
	logger Logger
	event_service.UnimplementedEventServiceServer
}

func NewEventService(app Application, logger Logger) *EventService {
	return &EventService{app: app, logger: logger}
}
