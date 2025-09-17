package event

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/divir112/otus_hw/internal/apperror"
	"github.com/divir112/otus_hw/internal/model"
)

func (s *EventServiceHTTP) GetEventListForDay(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	dataRequest, err := io.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			s.logger.Error("can't close request body: ", err)
		}
	}()
	if err != nil {
		s.logger.Error("can't read request body: ", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	requestModel := model.GetEventByDateRequest{}
	err = json.Unmarshal(dataRequest, &requestModel)
	if err != nil {
		s.logger.Error("can't unmarshal request body to GetEventByDateRequest model: ", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	events, err := s.app.GetEventsByDate(ctx, 1, requestModel.Date)
	if err != nil {
		if errors.Is(err, apperror.ErrIncorrectDate) {
			http.Error(w, "date is invalid", http.StatusBadRequest)
		}
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(events)
	if err != nil {
		s.logger.Error("can't marhal response body: ", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		s.logger.Error("can't write response get events list for day: ", err)
	}
}
