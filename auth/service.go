package auth

import (
	"github.com/dgrijalva/jwt-go"
)

// bagaimana cara membuat token (generate token) dan melakukan validasi token
type Service interface {
	// untuk generate token kita membutuhkan id, balikannya adalah string sebagai tokennya
	GenerateToken(userID int) (string, error)
}

// buat structnya
type jwtService struct {
}

// dengan menggunakkan ini kita bisa memanggil GenerateToken dari package manapun
func NewService() *jwtService {
	return &jwtService{}
}

// membuat secret key
var SECRET_KEY = []byte("CROWDFUNDING_s3cr3T_k3Y")

// function GenerateToken digunakan untuk membuat token JWT
func (s *jwtService) GenerateToken(userID int) (string, error) {
	// payload bisa juga disebut claim
	claim := jwt.MapClaims{}
	// dengan data key-nya adalah user_id dan value-nya adalah userID (parameter dari function generate token)
	claim["user_id"] = userID

	// generate tokennya
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// signature
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}
