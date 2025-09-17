package app

import (
	"context"
	"fmt"
	"time"

	"github.com/divir112/otus_hw/internal/apperror"
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
	GetEvent(ctx context.Context, id int) (model.Event, error)
	List(context.Context) ([]model.Event, error)
	CheckEventIsExists(ctx context.Context, startEvet, endEvent time.Time) (bool, error)
	GetEventsByDays(ctx context.Context, days int, date time.Time) ([]model.Event, error)
	GetReminderEvents(ctx context.Context) ([]model.Event, error)
}

func New(l Logger, s Storage) *App {
	return &App{
		logger:  l,
		storage: s,
	}
}

func (a *App) CreateEvent(ctx context.Context, event model.Event) (int, error) {
	if event.EndTime.Before(event.StartTime) {
		return 0, apperror.ErrTimeEndLessStart
	}

	exists, err := a.storage.CheckEventIsExists(ctx, event.StartTime, event.EndTime)
	if err != nil {
		return 0, fmt.Errorf("can't check existsing event: %w", err)
	}

	if exists {
		return 0, apperror.ErrDateBusy
	}

	id, err := a.storage.Add(ctx, event)
	if err != nil {
		return 0, fmt.Errorf("can't create event %w", err)
	}
	return id, nil
}

func (a *App) GetEvents(ctx context.Context) ([]model.Event, error) {
	events, err := a.storage.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get events %w", err)
	}
	return events, nil
}

func (a *App) UpdateEvents(ctx context.Context, id int, event model.Event) error {
	_, err := a.storage.GetEvent(ctx, id)
	if err != nil {
		return fmt.Errorf("cam't get event: %w", err)
	}

	err = a.storage.Update(ctx, id, event)
	if err != nil {
		return fmt.Errorf("can't update event %w", err)
	}

	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id int) error {
	err := a.storage.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("can't delete event %w", err)
	}

	return nil
}

func (a *App) GetEventsByDate(ctx context.Context, days int, date string) ([]model.Event, error) {
	currentDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("can't parse date: %w, %w", err, apperror.ErrIncorrectDate)
	}
	events, err := a.storage.GetEventsByDays(ctx, days, currentDate)
	if err != nil {
		return nil, fmt.Errorf("can't get events for %d days: %w", days, err)
	}

	return events, nil
}

// func (a *App) GetEventsShouldBeReminded(ctx context.Context) ([]model.Event, error) {
// 	events, err := a.storage.GetReminderEvents(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("can't ger reminder events: %w", err)
// 	}

// 	eventsShouldBeReminded := make([]model.Event, 0)
// 	dateNow :=
// 	for _, event := range events {
// 		pingDay := event.Date.Add(-(time.Hour * 24 * event.PingBefore))

// 		time.Parse(time.DateOnly, event.Date.Format()).
// 	}
// }

// TODO
