package event

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/divir112/otus_hw/internal/model"
)

func (s *EventServiceHTTP) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	splittedPath := strings.Split(r.URL.Path, "/events/")
	if len(splittedPath) < 1 {
		http.Error(w, `{"error": "id is required"}`, http.StatusBadRequest)
	}

	id, err := strconv.Atoi(splittedPath[1])
	if err != nil {
		s.logger.Error("can't atoi id: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	err = s.app.UpdateEvents(ctx, id, event)
	if err != nil {
		s.logger.Error("can't update event: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
