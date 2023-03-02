package domain

import "time"

type RefreshToken struct {
	ID            uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Refresh_Token string    `gorm:"size:255;not null;unique" json:"refresh_token"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}
