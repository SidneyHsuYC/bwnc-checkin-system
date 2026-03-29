package models

import "time"

// Class represents a class in the system
type Class struct {
	ID        int       `json:"id"`
	ClassName string    `json:"class_name"` // e.g., "Beginner Yoga", "Advanced Piano"
	StartDate time.Time `json:"start_date"`
	DayOfWeek string    `json:"day_of_week"` // Mon, Tue, Wed, Thu, Fri, Sat, Sun
	StartTime string    `json:"start_time"`  // e.g., "19:00"
	EndTime   string    `json:"end_time"`    // e.g., "21:00"
	StudentID *int      `json:"student_id"`  // Optional student leader (renamed from leader_id)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ClassWithLeader represents a class with leader details
type ClassWithLeader struct {
	ID        int       `json:"id"`
	ClassName string    `json:"class_name"`
	StartDate time.Time `json:"start_date"`
	DayOfWeek string    `json:"day_of_week"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	StudentID *int      `json:"student_id"` // Optional student leader (renamed from leader_id)
	Leader    *Student  `json:"leader,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
