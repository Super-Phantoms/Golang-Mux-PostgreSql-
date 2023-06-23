package domain

import (
	"strconv"

	"github.com/golangdevm/fullstack/dto"
)

const WITHDRAWAL = "withdrawal"

type Transaction struct {
	TransactionId   uint64  `gorm:"primary_key;auto_increment" json:"transaction_id" `
	AccountId       uint64  `gorm:"size:100;not null;" json:"account_id"`
	Amount          float64 `gorm:"size:255;not null;" json:"amount"`
	TransactionType string  `gorm:"size:100;not null;" json:"transaction_type"`
	TransactionDate string  `gorm:"size:255;not null;" json:"transaction_date"`
}

func (t Transaction) IsWithdrawal() bool {
	if t.TransactionType == WITHDRAWAL {
		return true
	}
	return false
}

func (t Transaction) ToDto() dto.TransactionResponse {
	return dto.TransactionResponse{
		TransactionId:   strconv.FormatUint(t.TransactionId, 10),
		AccountId:       strconv.FormatUint(t.AccountId, 10),
		Amount:          t.Amount,
		TransactionType: t.TransactionType,
		TransactionDate: t.TransactionDate,
	}
}
