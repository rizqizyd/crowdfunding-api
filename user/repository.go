package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	// contruct
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
}

// nilai db disini akan diisi sesuai dengan yang ada pada NewRepository
// struct repository harus memenuhi contruct interface repository
type repository struct {
	db *gorm.DB
}

// ketika NewRepository dipanggil, maka kita akan membuat object baru dari repository struct
func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// repository menyimpan ke database
func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

// function FindByEmail
func (r *repository) FindByEmail(email string) (User, error) {
	// mencari tabel user
	var user User

	err := r.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
