package model

import "time"

type Event struct {
	ID          int
	Title       string
	StartTime   time.Time
	EndTime     time.Time
	Description string
	UserID      int
}

type GetEventByDateRequest struct {
	Date string
}

type Notification struct {
	ID        int
	Title     string
	StartTime time.Time
	UserID    int
}
