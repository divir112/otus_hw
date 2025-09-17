package event

import (
	"context"

	"github.com/divir112/otus_hw/internal/model"
	desc "github.com/divir112/otus_hw/internal/pb/event_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *EventService) UpdateEvent(ctx context.Context, req *desc.Event) (*desc.UpdateEventsResponse, error) {
	id := req.Id
	if id == 0 {
		return nil, status.Error(codes.Internal, "id is required")
	}

	event := model.Event{
		Title:       req.Title,
		UserID:      int(req.UserId),
		Description: req.Description,
		StartTime:   req.StartTime.AsTime(),
		EndTime:     req.EndTime.AsTime(),
	}

	err := s.app.UpdateEvents(ctx, int(id), event)
	if err != nil {
		return nil, status.Error(codes.Internal, "can't update event")
	}

	return &desc.UpdateEventsResponse{}, nil
}
