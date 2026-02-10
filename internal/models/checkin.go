package models

import "time"

// Checkin represents a check-in record
type Checkin struct {
	ID          int       `json:"id"`
	StudentID   int       `json:"student_id"`
	EventID     int       `json:"event_id"`
	CheckedInAt time.Time `json:"checked_in_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// CheckinDetail represents a check-in with full student and event details
type CheckinDetail struct {
	ID          int       `json:"id"`
	Student     Student   `json:"student"`
	Event       Event     `json:"event"`
	CheckedInAt time.Time `json:"checked_in_at"`
}
