package model

import "time"

type Event struct {
	ID          int
	Header      string
	Date        time.Time
	DateEnd     time.Time `db:"dateend"`
	Description string
	Owner       string
}
