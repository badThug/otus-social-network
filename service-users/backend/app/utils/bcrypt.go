package utils

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(pwd, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
