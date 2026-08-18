package infrastructure

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService interface {
	Hash(password string) (string, error)
	ComparePassword(hashedPassword string, password string) error
}

type PasswordServiceImpl struct {
}

func NewPasswordService() *PasswordServiceImpl {
	return &PasswordServiceImpl{}
}

func (service *PasswordServiceImpl) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func (service *PasswordServiceImpl) ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
