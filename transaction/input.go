package transaction

import "api/user"

type GetCampaignTransactionsInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}

type CreateTransactionInput struct {
	Amount     int `json:"amount" binding:"required"`
	CampaignID int `json:"campaign_id" binding:"required"`

	// user yang melakukan request
	User user.User
}

/*
handler menerima data dari midtrans
bind data tersebut pakai context.ShouldBind ke dalam struct TransactionNotificationInput
struct TransactionNotificationInput akan jadi parameter dari service yang ada di dalam package payment
*/

// struct untuk menangkap notification dari midtrans
// 4 data penting yang dikirim oleh midtrans:
type TransactionNotificationInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"oerder_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
