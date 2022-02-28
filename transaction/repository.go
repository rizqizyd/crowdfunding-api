package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

// definisikan contruct/interface untuk menyimpan function"nya
type Repository interface {
	GetByCampaignID(campaignID int) ([]Transaction, error)
	GetByUserID(userID int) ([]Transaction, error)
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

// function GetByUserID (mencari data transactions yang dimiliki user tertentu)
func (r *repository) GetByUserID(userID int) ([]Transaction, error) {
	// supaya gorm tau kita akan mencari data di tabel transactions
	var transactions []Transaction

	// .Preload("Campaign.CampaignImages") => saat kita mengambil data transactions, data campaign yang berkaitan dengan suatu transaksi di load
	// sekaligus Campaign ini akan load CampaignImages (primary/gambar utama) yang berelasi dengan campaign tersebut
	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
