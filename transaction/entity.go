package transaction

import (
	"api/campaign"
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
	PaymentURL string
	User       user.User
	CreatedAt  time.Time
	updatedAt  time.Time

	// supaya transaction punya relasi dengan campaign lewat kolom campaign id
	Campaign campaign.Campaign
}
