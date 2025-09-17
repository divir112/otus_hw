package grpc

import (
	"context"
	"net"

	"github.com/divir112/otus_hw/internal/model"
	"github.com/divir112/otus_hw/internal/pb/event_service"
	"github.com/divir112/otus_hw/internal/server/grpc/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Logger interface {
	Info(string, ...any)
	Error(string, ...any)
}

type Application interface {
	CreateEvent(ctx context.Context, event model.Event) (int, error)
	GetEvents(ctx context.Context) ([]model.Event, error)
	UpdateEvents(ctx context.Context, id int, event model.Event) error
	DeleteEvent(ctx context.Context, id int) error
	GetEventsByDate(ctx context.Context, days int, date string) ([]model.Event, error)
}

type Server struct {
	app    Application
	logger Logger
	server *grpc.Server
}

func NewServer(app Application, logger Logger) *Server {
	return &Server{app: app, logger: logger}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", ":8085")
	if err != nil {
		return err
	}
	innterceptor := NewInterceptor(s.logger)
	server := grpc.NewServer(grpc.UnaryInterceptor(innterceptor.loggingMiddleware))
	reflection.Register(server)

	event_service.RegisterEventServiceServer(server, event.NewEventService(s.app, s.logger))
	err = server.Serve(lis)
	if err != nil {
		return err
	}
	s.server = server

	return nil
}

func (s *Server) Stop() {
	s.server.Stop()
}
