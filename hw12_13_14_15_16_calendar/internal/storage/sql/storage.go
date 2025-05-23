package sqlstorage

import (
	"context"
	"fmt"

	"github.com/divir112/otus_hw/internal/model"
	"github.com/georgysavva/scany/pgxscan"
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
		"INSERT INTO event (header, date, date_end, description, owner) VALUES ($1, $2, $3, $4, $5) RETURNING id",
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
		"UPDATE event SET header=$1, date=$2, date_end=$3, description=$4, owner=$5 where id=$6",
		event.Header, event.Date, event.DateEnd, event.Description, event.Owner, id,
	)
	if err != nil {
		return fmt.Errorf("can't update event %w", err)
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM event WHETE id=$1", id)
	if err != nil {
		return fmt.Errorf("can't delete event %w", err)
	}
	return nil
}

func (s *Storage) List(ctx context.Context) ([]model.Event, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, header, date, date_end, description, owner FROM event")
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
