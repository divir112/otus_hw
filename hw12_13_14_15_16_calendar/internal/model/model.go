package model

import "time"

type Event struct {
	ID          int
	Header      string
	Date        time.Time
	DateEnd     time.Time
	Description string
	Owner       string
}
