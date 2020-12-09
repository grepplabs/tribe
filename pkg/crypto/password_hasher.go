package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

const DefaultBCryptCost = 12

type PasswordHasher struct {
	bcryptCost int
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		bcryptCost: DefaultBCryptCost,
	}
}

func (h *PasswordHasher) WithBCryptCost(bcryptCost int) *PasswordHasher {
	h.bcryptCost = bcryptCost
	return h
}

func (h *PasswordHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.bcryptCost)
	return string(bytes), err
}

func (h *PasswordHasher) VerifyHashedPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
