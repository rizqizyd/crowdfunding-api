package transaction

import "time"

type CampaignTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

// function untuk single transaction
func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {
	// buat object baru
	formatter := CampaignTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt
	return formatter
}

// mengubah slice of Transactions menjadi slice of CampaignTransactionFormatter
func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionFormatter {
	// cek jika tidak ada transaction maka kembalikan array kosong
	if len(transactions) == 0 {
		return []CampaignTransactionFormatter{}
	}

	// jika ada datanya, maka lakukan perulangan
	var transactionFormatter []CampaignTransactionFormatter

	// dari satu data transaction kita bisa ubah ke dalam CampaignTransactionFormatter menggunakkan FormatCampaignTransaction
	for _, transaction := range transactions {
		formatter := FormatCampaignTransaction(transaction)
		transactionFormatter = append(transactionFormatter, formatter)
	}

	return transactionFormatter
}

// struct backer transaction
type UserTransactionFormatter struct {
	ID        int       `json:"id"`
	Amount    int       `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`

	// dependency campaign formatter
	Campaign CampaignFormatter `json:"campaign"`
}

type CampaignFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

// function backer transaction (single)
func FormatUserTransaction(transaction Transaction) UserTransactionFormatter {
	// buat formatter, instance of UserTransactionFormatter
	formatter := UserTransactionFormatter{}

	// mapping
	formatter.ID = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	// campaign formatter
	campaignFormatter := CampaignFormatter{}
	campaignFormatter.Name = transaction.Campaign.Name

	// default
	campaignFormatter.ImageURL = ""

	// cek apakah image punya gambar atau tidak
	if len(transaction.Campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	// set field campaign yang ada di dalam formatter
	formatter.Campaign = campaignFormatter

	return formatter
}

// function backer transaction (slice of transactions)
func FormatUserTransactions(transactions []Transaction) []UserTransactionFormatter {
	// cek jika tidak ada transaction maka kembalikan array kosong
	if len(transactions) == 0 {
		return []UserTransactionFormatter{}
	}

	// jika ada datanya, maka lakukan perulangan
	var transactionFormatter []UserTransactionFormatter

	// dari satu data transaction kita bisa ubah ke dalam UserTransactionFormatter menggunakkan FormatUserTransaction
	for _, transaction := range transactions {
		formatter := FormatUserTransaction(transaction)
		transactionFormatter = append(transactionFormatter, formatter)
	}

	return transactionFormatter
}

type TransactionFormatter struct {
	ID         int    `json:"id"`
	CampaignID int    `json:"campaign_id"`
	UserID     int    `json:"user_id"`
	Amount     int    `json:"amount"`
	Status     string `json:"status"`
	Code       string `json:"code"`
	PaymentURL string `json:"payment_url"`
}

// function untuk single transaction
func FormatTransaction(transaction Transaction) TransactionFormatter {
	// buat object baru
	formatter := TransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.CampaignID = transaction.CampaignID
	formatter.UserID = transaction.UserID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.Code = transaction.Code
	formatter.PaymentURL = transaction.PaymentURL
	return formatter
}
