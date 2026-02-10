package models

import "time"

// Event represents an event in the system
type Event struct {
	ID        int       `json:"id"`
	EventName string    `json:"event_name"`
	EventTime time.Time `json:"event_time"`
	EventType string    `json:"event_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
