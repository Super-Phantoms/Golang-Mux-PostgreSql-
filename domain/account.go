package domain

import (
	"strconv"

	"github.com/golangdevm/fullstack/dto"
	"github.com/golangdevm/fullstack/errs"
)

const dbTSLayout = "2006-01-02 15:04:05"

type Account struct {
	AccountId   uint64  `gorm:"primary_key;auto_increment" json:"account_id"`
	CustomerId  uint64  `gorm:"size:255;not null;" json:"customer_id"`
	OpeningDate string  `gorm:"size:255;not null;" json:"opening_date"`
	AccountType string  `gorm:"size:255;not null;" json:"account_type"`
	Amount      float64 `gorm:"size:255;not null;" json:"amount"`
	Status      string  `gorm:"size:255;not null;" json:"Status"`
}

func (a Account) ToNewAccountResponseDto() *dto.NewAccountResponse {
	return &dto.NewAccountResponse{strconv.FormatUint(a.AccountId, 10)}
}

//go:generate mockgen -destination=../mocks/domain/mockAccountRepository.go -package=domain github.com/golangdevm/fullstack/domain AccountRepository
type AccountRepository interface {
	Save(account Account) (*Account, *errs.AppError)
	SaveTransaction(transaction Transaction) (*Transaction, *errs.AppError)
	FindBy(accountId uint64) (*Account, *errs.AppError)
}

func (a Account) CanWithdraw(amount float64) bool {
	if a.Amount < amount {
		return false
	}
	return true
}

func NewAccount(customerId uint64, accountType string, amount float64) Account {
	return Account{
		CustomerId:  customerId,
		OpeningDate: dbTSLayout,
		AccountType: accountType,
		Amount:      amount,
		Status:      "1",
	}
}
