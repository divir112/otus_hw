package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/divir112/otus_hw/internal/apperror"
	"github.com/divir112/otus_hw/internal/model"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
		"INSERT INTO event (header, date, dateend, description, owner) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		event.Header, event.Date, event.DateEnd, event.Description, event.Owner,
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
		"UPDATE event SET header=$1, date=$2, dateend=$3, description=$4, owner=$5 where id=$6",
		event.Header, event.Date, event.DateEnd, event.Description, event.Owner, id,
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
	rows, err := s.pool.Query(ctx, "SELECT id, header, date, dateend, description, owner FROM event")
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
	row := s.pool.QueryRow(ctx, "SELECT id FROM event WHERE date <= $2 AND date_end >= $1")
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
	rows, err := s.pool.Query(ctx, "SELECT id FROM event WHERE date <= $2 AND date_end >= $1")
	if err != nil {
		return model.Event{}, fmt.Errorf("can't get events %w", err)
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
	rows, err := s.pool.Query(ctx, "SELECT id, header, date, dateend, description, owner FROM event WHERE date BETWEEN $1 AND $2", date, toDate)
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
