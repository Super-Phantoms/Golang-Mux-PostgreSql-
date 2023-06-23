package util

import "golang.org/x/crypto/bcrypt"

func Hash(password string) ([]byte, error) {
	// return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
