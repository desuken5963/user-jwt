package domain

import "time"

// User エンティティ
type User struct {
	ID        uint
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
