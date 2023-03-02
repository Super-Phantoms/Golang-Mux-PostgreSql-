package domain

import (
	"database/sql"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Login struct {
	Username   string         `gorm:"column:username"`
	CustomerId sql.NullString `gorm:"column:customer_id"`
	Accounts   sql.NullString `gorm:"column:accounts"`
	Role       string         `gorm:"column:role"`
}

func (l Login) ClaimsForAccessToken() AccessTokenClaims {
	if l.Accounts.Valid && l.CustomerId.Valid {
		return l.claimsForUser()
	} else {
		return l.claimsForAdmin()
	}
}

func (l Login) claimsForUser() AccessTokenClaims {
	accounts := strings.Split(l.Accounts.String, ",")
	return AccessTokenClaims{
		CustomerId: l.CustomerId.String,
		Accounts:   accounts,
		Username:   l.Username,
		Role:       l.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
		},
	}
}

func (l Login) claimsForAdmin() AccessTokenClaims {
	return AccessTokenClaims{
		Username: l.Username,
		Role:     l.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
		},
	}
}
