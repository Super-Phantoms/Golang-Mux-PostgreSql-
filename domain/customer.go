package domain

import (
	"strconv"

	"github.com/golangdevm/fullstack/dto"
	"github.com/golangdevm/fullstack/errs"
)

type Customer struct {
	ID          uint64 `gorm:"primary_key;auto_increment" json:"customer_id"`
	Name        string `gorm:"size:255;not null;unique" json:"name"`
	City        string `gorm:"size:255;not null;" json:"city"`
	Zipcode     string `json:"zipcode"`
	DateofBirth string `gorm:"size:255;not null;" json:"date_of_birth"`
	Status      string `gorm:"size:255;not null;" json:"status"`
}

func (c Customer) statusAsText() string {
	statusAsText := "active"
	if c.Status == "0" {
		statusAsText = "inactive"
	}
	return statusAsText
}

// strconv.FormatUint(c.ID, 10)
func (c Customer) ToDto() dto.CustomerResponse {
	return dto.CustomerResponse{
		Id:          strconv.FormatUint(c.ID, 10),
		Name:        c.Name,
		City:        c.City,
		Zipcode:     c.Zipcode,
		DateofBirth: c.DateofBirth,
		Status:      c.statusAsText(),
	}
}

type CustomerRepository interface {
	FindAll(status string) ([]Customer, *errs.AppError)
	ById(string) (*Customer, *errs.AppError)
}
