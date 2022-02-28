package transaction

import (
	"api/campaign"
	"errors"
)

type service struct {
	// dependency ke repository
	repository Repository
	// akses ke campaign repository
	campaignRepository campaign.Repository
}

// definisikan contruct/interface
// parameter GetCampaignTransactionsInput didapatkan dari struct pada input.go
type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionsByUserID(userID int) ([]Transaction, error)
}

// untuk instansiasi NewRepository pada repository.go
func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
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
