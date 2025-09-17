package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/divir112/otus_hw/internal/apperror"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"

	"github.com/divir112/otus_hw/internal/model"
	"github.com/georgysavva/scany/pgxscan"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) Add(ctx context.Context, event model.Event) (int, error) {
	row := s.pool.QueryRow(
		ctx,
		"INSERT INTO event (title, start_time, end_time, description, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		event.Title, event.StartTime, event.EndTime, event.Description, event.UserID,
	)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("can't create event %w", err)
	}

	return id, nil
}

func (s *Storage) Update(ctx context.Context, id int, event model.Event) error {
	_, err := s.pool.Exec(
		ctx,
		"UPDATE event SET title=$1, start_time=$2, end_time=$3, description=$4, user_id=$5 where id=$6",
		event.Title, event.StartTime, event.EndTime, event.Description, event.UserID, id,
	)
	if err != nil {
		return fmt.Errorf("can't update event %w", err)
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM event WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("can't delete event %w", err)
	}
	return nil
}

func (s *Storage) List(ctx context.Context) ([]model.Event, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, title, start_time, end_time, description, user_id FROM event")
	if err != nil {
		return nil, fmt.Errorf("can't get events %w", err)
	}
	var events []model.Event
	err = pgxscan.ScanAll(&events, rows)
	if err != nil {
		return nil, fmt.Errorf("can't scan events %w", err)
	}

	return events, nil
}

func (s *Storage) CheckEventIsExists(ctx context.Context, startEvet, endEvent time.Time) (bool, error) {
	row := s.pool.QueryRow(ctx, "SELECT id FROM event WHERE start_time <= $2 AND end_time >= $1", startEvet, endEvent)
	var id int
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("can't scan event %w", err)
	}
	return true, nil
}

func (s *Storage) GetEvent(ctx context.Context, id int) (model.Event, error) {
	rows, err := s.pool.Query(ctx, "SELECT id FROM event WHERE id = $1", id)
	if err != nil {
		return model.Event{}, fmt.Errorf("can't get event %w", err)
	}
	var event model.Event
	err = pgxscan.ScanOne(&event, rows)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Event{}, apperror.ErrNotFound
		}
		return model.Event{}, fmt.Errorf("can't scan event %w", err)
	}
	return event, nil
}

func (s *Storage) GetEventsByDays(ctx context.Context, days int, date time.Time) ([]model.Event, error) {
	toDate := date.Add(time.Duration(days) * (time.Hour * 24)).Format("2006-01-02")
	rows, err := s.pool.Query(ctx, "SELECT id, title, start_time, end_time, description, user_id FROM event WHERE start_time BETWEEN $1 AND $2", date, toDate)
	if err != nil {
		return nil, fmt.Errorf("can't get events by date %w", err)
	}
	events := make([]model.Event, 0)
	err = pgxscan.ScanAll(&events, rows)
	if err != nil {
		return nil, fmt.Errorf("can't scan events")
	}

	return events, nil
}

func (s *Storage) GetReminderEvents(ctx context.Context) ([]model.Event, error) {
	dateNow := time.Now()
	rows, err := s.pool.Query(ctx, "SELECT id, title, start_time, end_time, description, user_id FROM event WHERE (start_time - (ping_before * INTERVAL '1 day'))::date = $1::date", dateNow)
	if err != nil {
		return nil, fmt.Errorf("can't get events %w", err)
	}
	events := make([]model.Event, 0)
	err = pgxscan.ScanAll(&events, rows)
	if err != nil {
		return nil, fmt.Errorf("can't scan events")
	}

	return events, nil
}

func (s *Storage) DeleteOldEvents(ctx context.Context, expireTime time.Duration) error {
	expireDate := time.Now().Add(-expireTime)
	_, err := s.pool.Exec(ctx, "DELETE FROM event WHERE start_time < $1", expireDate)
	if err != nil {
		return fmt.Errorf("can't exec query: %w", err)
	}
	return nil
}
