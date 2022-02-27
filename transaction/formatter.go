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
