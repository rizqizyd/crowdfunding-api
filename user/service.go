package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// membuat contruct di service
type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
	GetUserByID(ID int) (User, error)
}

// struct service ini harus memenuhi contruct (interface service di atas)
// maka perlu bikin function dengan nama "Login" dengan parameternya "LoginInput" dan balikannya (User, error)
type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

// mapping struct input ke struct User
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

	// simpan struct User melalui repository
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

// input dari user di mapping ke struct CheckEmailInput
func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	// ambil nilai email yang ada di dalam struct CheckEmailInput, tampung dalam sebuah variabel "email"
	email := input.Email

	// masukkan email ke dalam repository
	// mencari email melalui repository berdasarkan email yang diinput oleh user
	// jika email user ditemukkan maka akan mengembalikan error
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	// jika email user tidak ditemukkan di database maka user bisa mendaftar menggunakan email tersebut
	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {
	// dapatkan user berdasarkan ID
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	// update attribute avatar file name
	// mengubah AvatarFileName sesuai dengan parameter fileLocation
	user.AvatarFileName = fileLocation

	// simpan perubahan avatar file name
	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

// ambil user dari db berdasarkan user_id lewat service
func (s *service) GetUserByID(ID int) (User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	// cek jika user tidak ditemukan
	if user.ID == 0 {
		// buat pesan error baru
		return user, errors.New("No user found on with that ID")
	}

	return user, nil
}
