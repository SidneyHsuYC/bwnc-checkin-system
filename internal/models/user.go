package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	LastName  string    `json:"last_name"`
	FirstName string    `json:"first_name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
