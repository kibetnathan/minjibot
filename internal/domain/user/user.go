package user

import "time"

type User struct {
	ID           int64
	UserID       string
	Email        string
	Passwordhash string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
