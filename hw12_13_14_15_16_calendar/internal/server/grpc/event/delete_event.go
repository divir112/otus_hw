package event

import (
	"context"

	desc "github.com/divir112/otus_hw/internal/pb/event_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *EventService) DeleteEvent(ctx context.Context, req *desc.DeleteEventRequest) (*desc.DeleteEventResponse, error) {
	id := req.Id

	err := s.app.DeleteEvent(ctx, int(id))
	if err != nil {
		s.logger.Error("error: %v", err)
		return nil, status.Error(codes.Internal, "can't delete event")
	}

	return &desc.DeleteEventResponse{}, nil
}
