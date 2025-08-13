package event

import (
	"context"

	desc "github.com/divir112/otus_hw/internal/pb/event_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *EventService) GetEvents(ctx context.Context, req *desc.GetEventsRequest) (*desc.GetEventsResponse, error) {
	events, err := s.app.GetEvents(ctx)

	if err != nil {
		return nil, status.Error(codes.Internal, "can't get errors")
	}

	var eventsResp []*desc.Event

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
