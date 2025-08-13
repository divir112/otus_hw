package event

import (
	"context"
	"time"

	"github.com/divir112/otus_hw/internal/model"
	desc "github.com/divir112/otus_hw/internal/pb/event_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *EventService) CreateEvent(ctx context.Context, req *desc.CreateEventRequest) (*desc.CreateEventResponse, error) {
	header := req.Header
	if header == "" {
		return nil, status.Error(codes.InvalidArgument, "header is required")
	}

	description := req.Description
	if description == "" {
		return nil, status.Error(codes.InvalidArgument, "description is required")
	}

	date := time.Now()

	dateEnd := req.DateEnd.AsTime()
	owner := req.Owner

	event := model.Event{
		Header:      header,
		Description: description,
		Date:        date,
		DateEnd:     dateEnd,
		Owner:       owner,
	}

	err := s.app.CreateEvent(ctx, event)

	if err != nil {
		s.logger.Error("error: %v", err)
		return nil, status.Error(codes.Internal, "can't create event")
	}

	return &desc.CreateEventResponse{}, nil
}
