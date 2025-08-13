package event

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/divir112/otus_hw/internal/model"
)

func (s *EventServiceHTTP) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	body := r.Body
	defer func() {
		err := r.Body.Close()
		if err != nil {
			s.logger.Error(fmt.Sprintf("can't close request body: %v", err))
		}
	}()

	dataRequest, err := io.ReadAll(body)
	if err != nil {
		s.logger.Error("can't read request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var event model.Event
	err = json.Unmarshal(dataRequest, &event)
	if err != nil {
		s.logger.Error("can't unmarshal request body to event model: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.app.CreateEvent(ctx, event)
	if err != nil {
		s.logger.Error("can't create event: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
