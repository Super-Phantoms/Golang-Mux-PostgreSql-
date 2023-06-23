package domain

import (
	"strconv"
	"time"

	"github.com/golangdevm/fullstack/dto"
)

type User struct {
	ID         uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Username   string    `gorm:"size:255;not null;unique" json:"username"`
	CustomerId uint32    `gorm:"size:100;not null;" json:"customer_id"`
	Email      string    `gorm:"size:100;not null;unique" json:"email"`
	Password   string    `gorm:"size:100;not null;" json:"password"`
	Role       string    `gorm:"size:100;not null;" json:"role"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u User) ToDto() dto.UserResponse {
	return dto.UserResponse{
		Id:         strconv.FormatUint(uint64(u.ID), 10),
		Username:   u.Username,
		CustomerId: strconv.FormatUint(uint64(u.CustomerId), 10),
		Email:      u.Email,
		Role:       u.Role,
	}
}
