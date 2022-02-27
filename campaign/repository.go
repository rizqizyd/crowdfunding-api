package campaign

import "gorm.io/gorm"

// definisikan terlebih dahulu contruct-nya
type Repository interface {
	// mengembalikan lebih dari 1 data campaign dari database
	FindAll() ([]Campaign, error)
	FindByUserID(userID int) ([]Campaign, error)
	FindByID(ID int) (Campaign, error)
	Save(campaign Campaign) (Campaign, error)
}

// definisikan struct
type repository struct {
	// field-nya db dengan tipe of *gorm.DB yang punya akses langusun ke database
	db *gorm.DB
}

// supaya struct bisa diakses dari luar package maka perlu dibuatkan instance
func NewRepository(db *gorm.DB) *repository {
	// return pembuatan instance baru, passing nilai db
	return &repository{db}
}

// selanjutnya bisa mulai implementasikan FindAll() dan FindByUserID() di struct repository
// menampilkan semua campaign yang tersedia
func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign
	// Preload akan me-load relasi dalam database, dalam kasus ini adalah "CampaignImages" (nama field)
	// dengan kondisi "campaign_images" (tabel pada database) filter is_primary = 1 untuk load satu gambar saja
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

// menampilkan semua campaign berdasarkan userID
func (r *repository) FindByUserID(userID int) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

// selanjutnya kita akan membuat instance dari struct repository melalui function NewRepository() yang kita panggil di dalam package main
// kemudian nanti akan kita tes apakah nilai yang dikembalikan function FindAll() dan FindUserByID sudah benar atau belum

// mengambil campaign berdasarkan id
func (r *repository) FindByID(ID int) (Campaign, error) {
	// variabel campaign yang merupakan tipe dari struct campaign
	var campaign Campaign

	err := r.db.Preload("User").Preload("CampaignImages").Where("id=?", ID).Find(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

// function create campaign
func (r *repository) Save(campaign Campaign) (Campaign, error) {
	// menyimpan campaign baru
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
