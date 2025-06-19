package memorystorage

import (
	"testing"
	"time"

	"github.com/divir112/otus_hw/internal/model"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage(t *testing.T) {
	event := model.Event{
		Header:      "test header",
		Date:        time.Now(),
		DateEnd:     time.Now().Add(time.Hour),
		Description: "test description",
		Owner:       "tester",
	}
	t.Run("test create event", func(t *testing.T) {
		events := make(map[int]model.Event)
		storage := New(events)
		id, err := storage.Add(nil, event) //nolint
		event.ID = 1
		require.NoError(t, err)
		require.Equal(t, 1, id)
		require.Equal(t, events[id], event)
	})

	t.Run("test delete event", func(t *testing.T) {
		events := make(map[int]model.Event)
		storage := New(events)
		id := 1
		events[1] = event

		err := storage.Delete(nil, id) //nolint
		require.NoError(t, err)
		_, exists := events[id]
		require.False(t, exists, "event does not delete")
	})

	t.Run("test update event", func(t *testing.T) {
		events := make(map[int]model.Event)
		storage := New(events)
		id := 1
		events[1] = event

		updatedEvent := model.Event{
			Header:      "epdated header",
			Date:        time.Now(),
			DateEnd:     time.Now().Add(time.Hour),
			Description: "epdated description",
			Owner:       "epdated",
		}
		err := storage.Update(nil, id, updatedEvent) //nolint
		require.NoError(t, err)
		require.Equal(t, updatedEvent, events[id])
	})

	t.Run("test get list events", func(t *testing.T) {
		events := make(map[int]model.Event)
		storage := New(events)
		events[1] = event

		eventList, err := storage.List(nil) //nolint
		require.NoError(t, err)
		require.Equal(t, []model.Event{event}, eventList)
	})
}
