package event

import (
	"context"
	"errors"

	"github.com/divir112/otus_hw/internal/apperror"
	desc "github.com/divir112/otus_hw/internal/pb/event_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *EventService) GetEventListForWeek(ctx context.Context, req *desc.GetEventListForDateRequest) (*desc.GetEventsResponse, error) {
	date := req.GetDate()
	if date == nil {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}
	strDate := date.AsTime().Format("2006-01-02")
	events, err := s.app.GetEventsByDate(ctx, 7, strDate)
	if err != nil {
		if errors.Is(err, apperror.ErrIncorrectDate) {
			return nil, status.Error(codes.InvalidArgument, "date is invalid")
		}
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	eventsResp := make([]*desc.Event, 0, len(events))
	for _, event := range events {
		eventsResp = append(eventsResp, &desc.Event{
			Id:          int64(event.ID),
			Header:      event.Header,
			Description: event.Description,
			Date:        timestamppb.New(event.Date),
			DateEnd:     timestamppb.New(event.DateEnd),
			Owner:       event.Owner,
		})
	}

	return &desc.GetEventsResponse{Events: eventsResp}, nil
}
