package campaign

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
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
