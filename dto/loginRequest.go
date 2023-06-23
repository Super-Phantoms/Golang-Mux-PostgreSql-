package dto

import (
	"errors"
	"html"
	"strings"

	"github.com/golangdevm/fullstack/util"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *LoginRequest) BeforeSave() error {
	hashedPassword, err := util.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *LoginRequest) Prepare() {
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
}

func (u *LoginRequest) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Username == "" {
			return errors.New("Required Username")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		return nil
	default:
		if u.Username == "" {
			return errors.New("Required Username")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		return nil
	}
}
