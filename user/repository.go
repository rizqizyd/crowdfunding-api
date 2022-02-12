package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	Save(user User) (User, error)
}

// nilai db disini akan diisi sesuai dengan yang ada pada NewRepository
type repository struct {
	db *gorm.DB
}

// ketika NewRepository dipanggil, maka kita akan membuat object baru dari repository struct
func NewRepository(db *gorm.DB) *repository  {
	return &repository{db}
}

// repository menyimpan ke database
func (r *repository) Save(user User) (User, error)  {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}