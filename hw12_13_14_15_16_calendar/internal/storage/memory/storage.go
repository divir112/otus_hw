package memorystorage

import (
	"context"
	"errors"
	"sync"

	"github.com/divir112/otus_hw/internal/model"
)

type Storage struct {
	mu     *sync.RWMutex
	events map[int]model.Event
	id     int
}

func New(events map[int]model.Event) *Storage {
	mu := &sync.RWMutex{}
	return &Storage{mu, events, 0}
}

func (s *Storage) Add(_ context.Context, event model.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.id++
	event.ID = s.id
	s.events[s.id] = event
	return s.id, nil
}

func (s *Storage) Update(_ context.Context, id int, event model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.events[id]
	if !ok {
		return errors.New("not found")
	}

	s.events[id] = event
	return nil
}

func (s *Storage) Delete(_ context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, id)
	return nil
}

func (s *Storage) List(_ context.Context) ([]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]model.Event, 0, len(s.events))
	for _, notification := range s.events {
		events = append(events, notification)
	}

	return events, nil
}
