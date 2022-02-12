package user

import (
	"time"
)

// pada framework lain, struct user ini dapat diibaraktan sebagai model
type User struct {
	ID int
	Name string
	Occupation string
	Email string
	PasswordHash string
	AvatarFileName string
	Role string
	CreatedAt time.Time
	UpdatedAt time.Time
}