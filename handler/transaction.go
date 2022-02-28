package handler

import (
	"api/helper"
	"api/transaction"
	"api/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	// dependency ke service
	service transaction.Service
}

// untuk instansiasi NewRepository pada service.go
func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

// function GetCampaignTransactions
func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	// tangkap input
	var input transaction.GetCampaignTransactionsInput

	// bind dengan uri yang ada di endpoint url
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get list of transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// dapatkan data currentUser
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	// panggil service
	transactions, err := h.service.GetTransactionsByCampaignID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get list of transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// balikan response berdasarkan format yang telah dibuat pada formatter
	response := helper.APIResponse("List of transactions", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

/*
Analisis campaign transaction:
-> parameter di uri
-> tangkap parameter mapping ke input struct
-> panggil service, input struct sebagai parameternya
-> service, berbekal campaign id bisa panggil repository
-> repository mencari data transaction suatu campaign
-> kembali ke handler untuk melakukkan formatting menggunakkan formatter
*/

// Analisis user transaction:
// GetUserTransactions
// handler
// ambil nilai user dari jwt/middleware
// service
// repository -> ambil data transactions (preload campaign)

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	// ambil data user yang melakukan request
	// ubah ke tipe (user.User)
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	// panggil service
	transactions, err := h.service.GetTransactionsByUserID(userID)
	if err != nil {
		response := helper.APIResponse("Failed to get list of backer transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of backer transactions", http.StatusOK, "success", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

/*
Analisis transaksi menggunakkan midtrans:
-> input jumlah donasi dari user (amount) - passing data user yang melakukan request
-> handler tangkap input, kemudian di mapping ke input struct
-> panggil service buat transaksi, panggil sistem midtrans untuk mendaftarkan transaksi (snapGateway.GetToken)
-> panggil repository create new transaction data
*/

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	// mapping inputan
	var input transaction.CreateTransactionInput

	// bind ke struct input
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create transaction", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// mengisi field user yang ada di struct CreateTransactionInput (input.go), diisi dengan user yang melakukan request
	// dapatkan data currentUser
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	// panggil CreateTransaction
	newTransaction, err := h.service.CreateTransaction(input)
	if err != nil {
		response := helper.APIResponse("Failed to create transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create transaction", http.StatusOK, "success", transaction.FormatTransaction(newTransaction))
	c.JSON(http.StatusOK, response)
}
