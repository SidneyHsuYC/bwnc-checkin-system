package models

import "time"

// Student represents a student account in the system
type Student struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	ClassInfo string    `json:"class_info,omitempty"` // Optional legacy field
	ClassID   *int      `json:"class_id,omitempty"`   // Optional class reference
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StudentSearchResult represents a student in search results
type StudentSearchResult struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	ClassInfo string `json:"class_info"`
	Email     string `json:"email"`
}
