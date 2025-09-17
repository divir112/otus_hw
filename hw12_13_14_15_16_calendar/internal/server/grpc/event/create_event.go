package event

import (
	"context"
	"errors"

	"github.com/divir112/otus_hw/internal/apperror"
	"github.com/divir112/otus_hw/internal/model"
	desc "github.com/divir112/otus_hw/internal/pb/event_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *EventService) CreateEvent(ctx context.Context, req *desc.CreateEventRequest) (*desc.CreateEventResponse, error) {
	title := req.Title
	if title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	description := req.Description
	if description == "" {
		return nil, status.Error(codes.InvalidArgument, "description is required")
	}

	date := req.StartTime.AsTime()

	dateEnd := req.EndTime.AsTime()
	userID := req.UserId

	event := model.Event{
		Title:       title,
		Description: description,
		StartTime:   date,
		EndTime:     dateEnd,
		UserID:      int(userID),
	}

	id, err := s.app.CreateEvent(ctx, event)

	if err != nil {
		s.logger.Error("error: %v", err)
		if errors.Is(err, apperror.ErrTimeEndLessStart) {
			return nil, status.Error(codes.AlreadyExists, "date is busy")
		}
		return nil, status.Error(codes.Internal, "can't create event")
	}

	return &desc.CreateEventResponse{Id: int64(id)}, nil
}
