package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)

	// parameter yang digunakan untuk Update campaign adalah struct pada input.go
	// parameter inputID merupakan tipe dari GetCampaignDetailInput, dst.
	UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error)

	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
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

// implementasi dari contruct Update
func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error) {
	// ambil data campaign berdasarkan id (nilai yang masih lama)
	campaign, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	}

	// passing current user dilakukan supaya kita bisa cek bahwa user yang melakukan update adalah user yang memiliki campaign tersebut
	// jika user yang melakukan request bukan user yang memiliki data maka kita kembalikan error
	if campaign.UserID != inputData.User.ID {
		return campaign, errors.New("Not an owner of the campaign")
	}

	// tangkap parameternya kemudian mapping ke object campaign yang ingin di update
	// nilai data baru
	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	// simpan ke dalam database menggunakan repository yang telah dibuat
	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}

// implementasi SaveCampaignImage
func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error) {
	// ambil data campaign berdasarkan id
	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return CampaignImage{}, err
	}

	// jika yang melakukan request bukan user yang bikin campaign, maka dia tidak bisa upload
	if campaign.UserID != input.User.ID {
		return CampaignImage{}, errors.New("Not an owner of the campaign")
	}

	isPrimary := 0
	// cek apakah perlu mengubah is_primary jadi false atau tidak
	if input.IsPrimary {
		isPrimary = 1
		// jika primary true diubah menjadi false
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return CampaignImage{}, err
		}
	}

	// save campaign image baru
	campaignImage := CampaignImage{}
	// proses mapping
	campaignImage.CampaignID = input.CampaignID
	campaignImage.IsPrimary = isPrimary
	campaignImage.FileName = fileLocation

	// panggil repository
	newCampaignImage, err := s.repository.CreateImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}
