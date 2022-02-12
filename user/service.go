package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// membuat contruct di service
type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
}

// struct service ini harus memenuhi contruct (interface service di atas)
// maka perlu bikin function dengan nama "Login" dengan parameternya "LoginInput" dan balikannya (User, error)
type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

// function Login
func (s *service) Login(input LoginInput) (User, error) {
	// ambil nilai email dan password
	email := input.Email
	password := input.Password

	// langkah selanjutnya adalah mencari user dengan alamat email tersebut
	// memanfaatkan repository yang sudah dibuat sebelumnya
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	// cek jika user ada nilainya atau tidak
	if user.ID == 0 {
		// buat pesan error baru
		return user, errors.New("No user found on that email")
	}

	// mencocokkan password yang dimasukkan oleh user dengan passowrd yang ada di database
	// menggunakkan fungsi bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	// cek password match atau tidak
	if err != nil {
		return user, err
	}

	// jika password match maka dibalikkan tanpa err
	return user, nil
}

// mapping struct input ke struct User
// simpan struct User melalui repository
