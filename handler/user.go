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

func (h *userHandler) RegisterUser(c *gin.Context) {
	// tangkap input dari user
	// map input dari user ke struct RegisterUserInput
	// struct di atas kita passing sebagai parameter service

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
