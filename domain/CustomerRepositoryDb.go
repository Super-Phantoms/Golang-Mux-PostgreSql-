package domain

import (
	"database/sql"

	"github.com/golangdevm/fullstack/errs"
	"github.com/golangdevm/fullstack/logger"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

type CustomerRepositoryDb struct {
	client *gorm.DB
}

func (d CustomerRepositoryDb) FindAll(status string) ([]Customer, *errs.AppError) {
	var err error
	customers := []Customer{}
	// customers := make([]Customer, 0)

	if status == "" {
		err = d.client.Debug().Model(&Customer{}).Limit(100).Find(&customers).Error
		// findAllSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers"
		// err = d.client.Select(&customers, findAllSql)
	} else {
		err = d.client.Debug().Model(&Customer{}).Where("status = ?", status).Take(&customers).Error
		// findAllSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers where status = ?"
		// err = d.client.Select(&customers, findAllSql, status)
	}

	if err != nil {
		logger.Error("Error while querying customers table " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return customers, nil
}

func (d CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {
	// customerSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers where customer_id = ?"

	var c Customer
	err := d.client.Debug().Model(&Customer{}).Where("id = ?", id).Take(&c).Error
	// err := d.client.Get(&c, customerSql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer not found")
		} else {
			logger.Error("Error while scanning customer " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &c, nil
}

func NewCustomerRepositoryDb(dbClient *gorm.DB) CustomerRepositoryDb {
	return CustomerRepositoryDb{dbClient}
}
