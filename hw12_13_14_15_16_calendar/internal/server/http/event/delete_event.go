package event

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

func (s *EventServiceHTTP) DeleteEvent(w http.ResponseWriter, r *http.Request) {
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

	err = s.app.DeleteEvent(ctx, id)
	if err != nil {
		s.logger.Error("can't delete event: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
