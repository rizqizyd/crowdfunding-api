package payment

import (
	"api/user"
	"strconv"

	midtrans "github.com/veritrans/go-midtrans"
)

// documentation: https://github.com/veritrans/go-midtrans
type service struct {
}

type Service interface {
	// object transaksi untuk melakukan request, passing data user
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

func NewService() *service {
	return &service{}
}

// implementasi GetPaymentURL
func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	// daftarkan transaksinya
	midclient := midtrans.NewClient()
	midclient.ServerKey = "YOUR-VT-SERVER-KEY"
	midclient.ClientKey = "YOUR-VT-CLIENT-KEY"
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	// saat akan melakukan request, sebelum mendapatkan tokennya perlu memasukkan object request (snapReq)
	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	// GetToken untuk mendapatkan redirect url
	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}
