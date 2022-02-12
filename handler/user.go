package handler

import (
	"api/auth"
	"api/helper"
	"api/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// field disini dibuat supaya setiap service dapat dipakai di dalam userHandler
type userHandler struct {
	userService user.Service
	authService auth.Service
}

// userService akan di passing menjadi userService yang ada di struct
func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
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

	// panggil authService
	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// memanggil formatter
	formatter := user.FormatUser(newUser, token)

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

	// panggil authService
	token, err := h.authService.GenerateToken(loggedinUser.ID)
	if err != nil {
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// memanggil formatter. balikkan loggedinUser ke dalam format json. user di format dalam bentuk user formatter
	formatter := user.FormatUser(loggedinUser, token)

	// memanggil response pada helper. passing ke dalam api response
	response := helper.APIResponse("Successfuly loggedin", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

// cek email API - apakah sudah terdaftar atau belum
func (h *userHandler) CheckEmailAvaliability(c *gin.Context) {
	/*
		Steps:
		1. cek apakah ada input email dari user
		2. input email di mapping ke struct input (prosesnya ada di handler. ada berbagai macam validasi)
		3. struct input di passing ke service
		4. service akan memanggil repository untuk menentukkan apakah email sudah ada atau belum
		5. repository melakukan query ke database
	*/

	// menangkap input yang dimasukkan user
	var input user.CheckEmailInput

	// proses mapping/binding
	err := c.ShouldBindJSON(&input)
	// validasi error
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		// status StatusUnprocessableEntity: 422
		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// memanggil service yang memiliki 2 balikan yaitu boolean dan error
	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// mengembalikan response
	data := gin.H{
		"is_available": isEmailAvailable,
	}

	// mengecek isEmailAvailable true or false untuk menentukan message yang ada di meta
	metaMessage := "Email has been registered"
	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

// function upload avatar
func (h *userHandler) UploadAvatar(c *gin.Context) {
	// buat 2 function di repository
	// repository ambil data user yang ID = x
	// repository update data user simpan lokasi file

	// menangkap input dari user (parameternya bukan json tapi form body (string))
	file, err := c.FormFile("avatar")
	// balikan response
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// simpan gambarnya di folder "images/" + fileName
	userID := 1 // dapet dari JWT nanti
	// path := "images/" + file.Filename // nama file tanpa id (bisa konflik dengan user lain yang mengupload nama file sama)
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename) // nama file dengan id

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// di service kita panggil repository
	// JWT (sementara hardcode, seakan2 user yg login ID = x)
	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar successfuly uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}
