package bcrypt

import (
	"github.com/grepplabs/tribe/pkg/crypto"
	"golang.org/x/crypto/bcrypt"
)

const DefaultBCryptCost = 10

type bcryptPasswordHasher struct {
	bcryptCost int
}

func NewPasswordHasher(options ...Option) (crypto.PasswordHasher, error) {
	hasher := &bcryptPasswordHasher{
		bcryptCost: DefaultBCryptCost,
	}
	for _, option := range options {
		if err := option(hasher); err != nil {
			return nil, err
		}
	}
	return hasher, nil
}

func (h *bcryptPasswordHasher) WithBCryptCost(bcryptCost int) crypto.PasswordHasher {
	h.bcryptCost = bcryptCost
	return h
}

func (h *bcryptPasswordHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.bcryptCost)
	return string(bytes), err
}

func (h *bcryptPasswordHasher) VerifyHashedPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
