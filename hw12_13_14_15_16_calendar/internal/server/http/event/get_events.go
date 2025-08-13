package event

import (
	"context"
	"encoding/json"
	"net/http"
)

func (s *EventServiceHTTP) GetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	events, err := s.app.GetEvents(ctx)

	if err != nil {
		s.logger.Error("can't get events: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	responseBytes, err := json.Marshal(events)
	if err != nil {
		s.logger.Error("can't marshal response get events: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)

	if err != nil {
		s.logger.Error("can't write response get events: ", err)
	}
}
