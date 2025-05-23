package app

import (
	"context"
	"fmt"

	"github.com/divir112/otus_hw/internal/model"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(string, ...any)
}

type Storage interface {
	Add(context.Context, model.Event) (int, error)
	Update(context.Context, int, model.Event) error
	Delete(context.Context, int) error
	List(context.Context) ([]model.Event, error)
}

func New(l Logger, s Storage) *App {
	return &App{
		logger:  l,
		storage: s,
	}
}

func (a *App) CreateEvent(ctx context.Context, id int, header string) error {
	_, err := a.storage.Add(ctx, model.Event{ID: id, Header: header})
	if err != nil {
		return fmt.Errorf("can't create event %w", err)
	}
	return nil
}

// TODO
