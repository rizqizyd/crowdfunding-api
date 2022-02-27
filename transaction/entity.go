package transaction

import (
	"api/user"
	"time"
)

// struct ini adalah representasi dari tabel transaction yang ada di database
type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	User       user.User
	CreatedAt  time.Time
	updatedAt  time.Time
}
