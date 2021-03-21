package bcrypt

import (
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestPasswordHasher(t *testing.T) {
	tests := []struct {
		name     string
		password string
		cost     int
		err      error
	}{
		{name: "empty", password: ""},
		{name: "1 char", password: mustGenerateRandomString(1)},
		{name: "8 chars", password: mustGenerateRandomString(8)},
		{name: "64 chars", password: mustGenerateRandomString(64)},
		{name: "8 chars cost 10", password: mustGenerateRandomString(8), cost: 10},
		{name: "64 chars cost 10", password: mustGenerateRandomString(64), cost: 10},
		{name: "255 chars cost 10", password: mustGenerateRandomString(255), cost: 10},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a := assert.New(t)

			cost := tc.cost
			if cost <= 0 {
				cost = DefaultBCryptCost
			}
			hasher, err := NewPasswordHasher(WithBCryptCost(tc.cost))
			a.Nil(err)
			a.Equal(tc.cost, hasher.(*bcryptPasswordHasher).bcryptCost)

			passwordHash, err := hasher.HashPassword(tc.password)
			a.Equal(tc.err, err)
			if err == nil {
				a.True(hasher.VerifyHashedPassword(tc.password, passwordHash), "Verify hash failed in test '%s'", tc.name)
			}
		})
	}
}

func mustGenerateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret)
}
