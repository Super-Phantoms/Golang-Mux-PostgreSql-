package domain

import (
	"github.com/golangdevm/fullstack/errs"
	"github.com/golangdevm/fullstack/logger"
	"github.com/jinzhu/gorm"
)

type AccountRepositoryDb struct {
	client *gorm.DB
}

func (d AccountRepositoryDb) Save(a Account) (*Account, *errs.AppError) {

	// saveErr := req.BeforeSave()
	// if saveErr != nil {
	// 	logger.Error("Error while hashing password: " + saveErr.Error())
	// 	return &User{}, errs.NewValidationError("hasing error")
	// }

	// req.Prepare()

	// valErr := req.Validate("register")
	// if valErr != nil {
	// 	logger.Error("validation error: " + valErr.Error())
	// 	return &User{}, errs.NewValidationError("validation error")
	// }

	newAccount := Account{
		CustomerId:  a.CustomerId,
		OpeningDate: a.OpeningDate,
		Amount:      a.Amount,
		Status:      a.Status,
	}

	err := d.client.Debug().Model(&Account{}).Create(&newAccount).Error
	if err != nil {

		logger.Error("Error while creating new account: " + err.Error())
		return &Account{}, errs.NewUnexpectedError("Unexpected error from database")
	} else {
		return &newAccount, nil
	}
}

/**
 * transaction = make an entry in the transaction table + update the balance in the accounts table
 */
func (d AccountRepositoryDb) SaveTransaction(t Transaction) (*Transaction, *errs.AppError) {

	tx := d.client.Begin()
	insertSql := `INSERT INTO transactions (account_id, amount, transaction_type, transaction_date) values (?, ?, ?, ?)`
	insertErr := tx.Exec(
		insertSql, t.AccountId, t.Amount, t.TransactionType, t.TransactionDate,
	).Debug().Error
	if insertErr != nil {
		logger.Error("Error while inserting transaction " + insertErr.Error())
		return nil, errs.NewUnexpectedError("Unexpected transaction error")
	}

	withdrawalSql := `UPDATE accounts SET amount = amount - ? where account_id = ?`
	depositSql := `UPDATE accounts SET amount = amount + ? where account_id = ?`
	// updating account balance
	var err error
	if t.IsWithdrawal() {
		err = tx.Exec(withdrawalSql, t.Amount, t.AccountId).Debug().Error
	} else {
		err = tx.Exec(depositSql, t.Amount, t.AccountId).Debug().Error
	}

	// in case of error Rollback, and changes from both the tables will be reverted
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	// commit the transaction when all is good
	commitErr := tx.Commit().Debug().Error

	if commitErr != nil {
		tx.Rollback()
		logger.Error("Error while commiting transaction for bank account: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	// var tran Transaction
	// tranErr := result.Debug().Model(&Transaction{}).Find(&t).Error
	// if tranErr != nil {
	// 	logger.Error("Error while querying transaction table " + tranErr.Error())
	// 	return nil, errs.NewUnexpectedError("Unexpected database error")
	// }

	// Getting the latest account information from the accounts table
	account, appErr := d.FindBy(t.AccountId)
	if appErr != nil {
		return nil, appErr
	}

	// updating the transaction struct with the latest balance
	t.Amount = account.Amount
	return &t, nil
}

func (d AccountRepositoryDb) FindBy(accountId uint64) (*Account, *errs.AppError) {
	// sqlGetAccount := "SELECT account_id, customer_id, opening_date, account_type, amount from accounts where account_id = ?"
	var account Account
	err := d.client.Debug().Model(&Account{}).Where("account_id = ?", accountId).Take(&account).Error
	// err := d.client.Get(&account, sqlGetAccount, accountId)
	if err != nil {
		logger.Error("Error while fetching account information: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return &account, nil
}

func NewAccountRepositoryDb(dbClient *gorm.DB) AccountRepositoryDb {
	return AccountRepositoryDb{dbClient}
}
