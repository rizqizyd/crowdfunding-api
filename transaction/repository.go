package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

// definisikan contruct/interface untuk menyimpan function"nya
type Repository interface {
	GetByCampaignID(campaignID int) ([]Transaction, error)
}

// untuk instansiasi NewRepository pada main.go
func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// function GetByCampaignID
func (r *repository) GetByCampaignID(campaignID int) ([]Transaction, error) {
	var transaction []Transaction

	// urutkan dengan menggunakkan function .Order berdasarkan id yang paling besar (donasi terakhir)
	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("id desc").Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
