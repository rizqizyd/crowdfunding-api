package handler

import (
	"api/helper"
	"api/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

// userService akan di passing menjadi userServiec yang ada di struct
func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

// user handler register
func (h *userHandler) RegisterUser(c *gin.Context) {
	/*
		Steps:
		1. tangkap input dari user
		2. map input dari user ke struct RegisterUserInput
		3. struct di atas kita passing sebagai parameter service

		*untuk membuat sebuah endpoint baru, kita akan mulai dari bagian yang terbawah (layer terbawah) yaitu:
		repository, (input dan service), handler, formatter
	*/

	// object input akan di mapping yang tadinya json menjadi RegisterUserInput
	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		// menangani validasi error, membuat array dan menambah data array (error) melalui perulangan
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// memasukkan input ke dalam service
	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// memanggil formatter
	formatter := user.FormatUser(newUser, "tokentokentoken")

	// memanggil response pada helper
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

// user handler login
func (h *userHandler) Login(c *gin.Context) {
	/*
		Steps:
		1. user memasukkan input (email dan password)
		2. input ditangkap handler
		3. mapping dari input user ke input struct
		4. input struct di passing ke service
		5. di service kita akan mencari dengan bantuan repository user dengan email x
		6. mencocokkan password
	*/

	// langkah terakhir: menangkap input dari user, kemudian kita map/bind ke dalam struct LoginInput
	var input user.LoginInput
	// proses mapping/binding
	err := c.ShouldBindJSON(&input)
	// validasi error
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		// status StatusUnprocessableEntity: 422
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// menampung user yang login
	loggedinUser, err := h.userService.Login(input)
	// cek jika user tidak ditemukkan atau salah password
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// memanggil formatter. balikkan loggedinUser ke dalam format json. user di format dalam bentuk user formatter
	formatter := user.FormatUser(loggedinUser, "tokentokentoken")

	// memanggil response pada helper. passing ke dalam api response
	response := helper.APIResponse("Successfuly loggedin", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}
