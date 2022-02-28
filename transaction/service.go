package transaction

import (
	"api/campaign"
	"api/payment"
	"errors"
	"strconv"
)

type service struct {
	// dependency ke repository
	repository Repository
	// akses ke campaign repository
	campaignRepository campaign.Repository
	// payment service
	paymentService payment.Service
}

// definisikan contruct/interface
// parameter GetCampaignTransactionsInput didapatkan dari struct pada input.go
type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionsByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
	ProcessPayment(input TransactionNotificationInput) error
}

// untuk instansiasi NewRepository pada repository.go
func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
}

// function GetTransactionsByCampaignID
func (s *service) GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {
	// authorization
	// get campaign -> check campaign.userid != user id yang melakukan request
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	// jika user bukan pemilik campaign, maka dia tidak bisa melihat data transaction
	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("Not an owner of the campaign")
	}

	// hanya memanggil si repository dengan function-nya GetByCampaignID
	transaction, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

// function GetTransactionsByUserID
func (s *service) GetTransactionsByUserID(userID int) ([]Transaction, error) {
	// panggil repository
	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

// function CreateTransaction
func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	// buat object transaction
	transaction := Transaction{}
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID

	// secara default user telah melakukan transaksi namun belum dibayar
	transaction.Status = "pending"
	// jika ingin membuat kode transaksi yang unique, untuk saat ini di kosongkan saja
	// transaction.Code = ""

	// panggil repository
	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	// mapping dari transaction yang ada di transaction, menjadi transaction yang ada di payment
	paymenTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	// panggil GetPaymentURL (hubungi midtrans untuk mendapatkan paymentURL)
	paymentURL, err := s.paymentService.GetPaymentURL(paymenTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}

	// dapatkan paymentURL -> update data transaction supaya punya data paymentURL yang didapatkan dari midtrans
	// simpan paymentURL ke dalam object newTransaction
	newTransaction.PaymentURL = paymentURL

	// panggil repository
	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}

// implementasi function ProcessPayment
func (s *service) ProcessPayment(input TransactionNotificationInput) error {
	// transaction_id digunakan untuk mengambil data transaction dengan id yang bersangkutan
	// ubah dari string ke integer
	transaction_id, _ := strconv.Atoi(input.OrderID)

	// mendapatkan transaction, id nya didapatkan dari inputan yang dikirim oleh midtrans
	transaction, err := s.repository.GetByID(transaction_id)
	if err != nil {
		return err
	}

	// update status transaksi pada database. kondisinya dikirim dari notifikasi midtrans
	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	// update data transaksi
	updatedTransaction, err := s.repository.Update(transaction)
	if err != nil {
		return err
	}

	// Update data campaign
	// ambil data campaign
	campaign, err := s.campaignRepository.FindByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}

	// jika ada transaksi yang berubah jadi paid, perlu kita update data backer count nya
	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		// panggil repository campaign, passing data campaign
		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	// jika semua berjalan lancar dan tidak ada error
	return nil
}
