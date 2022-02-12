package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// bagaimana cara membuat token (generate token) dan melakukan validasi token
type Service interface {
	// untuk generate token kita membutuhkan id, balikannya adalah string sebagai tokennya
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
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
	// HS256 adalah salah satu bentuk dari HMAC
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// signature
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

// validasi token
func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	// parse token
	// parameter: encodedToken (token yang panjang), function yang mengembalikan interface dan error
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		// cek apakah token tersebut tipenya adalah HMAC / apakah sama metode yang digunaka
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		// apakah benar token-nya dibuat dengan SECRET_KEY yang kita punya?
		if !ok {
			return nil, errors.New("Invalid token")
		}

		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return token, err
	}

	// kembalikan token jika berhasil di validasi
	return token, nil
}
