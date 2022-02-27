package handler

import (
	"api/campaign"
	"api/helper"
	"api/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
Analisis langkah-langkah membuat endpoint campaign:
-> tangkap parameter di handler
-> handler ke service
-> service yang menentukan repository (method) mana yang di call
-> repository (function/method) : FindAll, FindByUserID
-> db
*/

type campaignHandler struct {
	// kita membutuhkan bantuan service, maka kita bikin field namanya service
	// dimana dia adalah tipe dari package campaign dan tipenya adalah interface Service
	service campaign.Service
}

// untuk membuat object/struct dari campaignHandler ini yang nantinya akan dipanggil di dalam main.go
func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

// function yang akan di rout ke api/v1/campaigns
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	// tangkap parameter di handler (user id)
	// convert dari string ke integer
	userID, _ := strconv.Atoi(c.Query("user_id"))

	// panggil service GetCampaigns
	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// campaigns yang didapatkan kita ubah menjadi array/slice of campaign formatter menggunakkan fungsi FormatCampaigns
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

// function GetCampaign untuk mendapatkan detail campaign (api/v1/campaigns/id)
func (h *campaignHandler) GetCampaign(c *gin.Context) {
	/*
		-> api/v1/campaigns/:id
		-> handler: mapping id pada url ke struct input untuk dimasukkan ke service, call formatter
		-> service: inputnya struct input untuk mengangkap id pada url, memanggil repo
		-> repository: get campaign by id
	*/

	// variabel input akan menyimpan id campaign
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// panggil service
	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// balikan response berdasarkan format yang telah dibuat pada formatter
	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

/*
Analisis langkah-langkah create campaign
-> tangkap parameter dari user ke input struct
-> ambil current user dari jwt/handler (untuk mengetahui user pembuat campaign)
-> panggil service, parameternya adalah input struct yang telah di mapping
	- buat slug (otomatis berdasarkan nama campaign)
-> panggil repository untuk simpan data campaign baru
*/

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	// variabel input digunakan untuk menangkap parameter yang dikirim oleh user
	var input campaign.CreateCampaignInput

	// lakukan pengecekan apakah ada error saat melakukan ShouldBindJSON
	err := c.ShouldBindJSON(&input)
	if err != nil {
		// menangani validasi error, membuat array dan menambah data array (error) melalui perulangan
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// dapatkan data currentUser
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	// sampai sini berarti field yang ada di input sudah lengkap

	// panggil service
	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create campaign", http.StatusOK, "success", campaign.FormatCampaign((newCampaign)))
	c.JSON(http.StatusOK, response)
}
