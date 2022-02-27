package campaign

import (
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
}

// struct repository memiliki denpendency (ketergantungan) field terhadap
// interface Repository yang ada di dalam package campaign
type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

// implementasi dari contruct GetCampaigns
// service GetCampaigns ini nantinya akan dipanggil oleh handler campaign.go yang data userID nya di dapatkan melalui query parameter "user_id"
func (s *service) GetCampaigns(userID int) ([]Campaign, error) {
	// jika userID != 0, maka kita hanya mengambil data campaign berdasarkan id yang bersangkutan
	if userID != 0 {
		campaigns, err := s.repository.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}

	// jika userID = 0, maka semua campaign akan ditampilkan
	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

// service GetCampaignByID untuk mendapatkan data campaign berdasarkan id
func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	// memanggil repository
	campaign, err := s.repository.FindByID(input.ID)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

// buat campaign
// dari inputan user kita mapping ke CreateCampaignInput
func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	// kemudian kita mapping lagi dari CreateCampaignInput menjadi object Campaign
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount

	// mapping User (hanya membutuhkan id nya)
	campaign.UserID = input.User.ID

	// pembuatan slug (menggunakan library slug)
	// menggabungkan nama campaign dengan user id
	// Nama Campaign id=10 => nama-campaign-10
	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugCandidate)

	// panggil repository
	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}
