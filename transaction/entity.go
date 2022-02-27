package transaction

import "time"

// struct ini adalah representasi dari tabel transaction yang ada di database
type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	CreatedAt  time.Time
	updatedAt  time.Time
}
