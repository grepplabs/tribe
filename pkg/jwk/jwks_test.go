package jwk

import (
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestJwktGenerate(t *testing.T) {
	gen := NewJWKSGenerator()

	tests := []struct {
		id  string
		alg string
		use string
	}{
		{alg: "HS256", use: "sig"},
		{alg: "HS384", use: "sig"},
		{alg: "HS512", use: "sig"},
		{alg: "RS256", use: "sig"},
		{alg: "RS384", use: "sig"},
		{alg: "RS512", use: "sig"},
		{alg: "ES256", use: "sig"},
		{alg: "ES384", use: "sig"},
		{alg: "ES512", use: "sig"},
		{alg: "PS256", use: "sig"},
		{alg: "PS384", use: "sig"},
		{alg: "PS512", use: "sig"},
		{alg: "HS256", use: "enc"},
		{alg: "HS384", use: "enc"},
		{alg: "HS512", use: "enc"},
		{alg: "RS256", use: "enc"},
		{alg: "RS384", use: "enc"},
		{alg: "RS512", use: "enc"},
		{alg: "ES256", use: "enc"},
		{alg: "ES384", use: "enc"},
		{alg: "ES512", use: "enc"},
		{alg: "PS256", use: "enc"},
		{alg: "PS384", use: "enc"},
		{alg: "PS512", use: "enc"},

		{alg: "HS256", use: "sig", id: uuid.NewString()},
		{alg: "ES256", use: "sig", id: uuid.NewString()},
		{alg: "RS256", use: "sig", id: uuid.NewString()},
		{alg: "PS256", use: "sig", id: uuid.NewString()},
	}
	for _, tc := range tests {
		name := fmt.Sprintf("%s-%s", tc.alg, tc.use)
		if tc.id != "" {
			name = fmt.Sprintf("%s-%s", name, tc.id)
		}
		t.Run(name, func(t *testing.T) {
			a := assert.New(t)
			jwksset, err := gen.Generate(tc.id, tc.alg, tc.use)
			a.Nil(err)
			a.False(len(jwksset.Keys) == 0)
			if strings.HasPrefix(tc.alg, "HS") {
				a.Equal(1, len(jwksset.Keys))
			} else {
				a.Equal(2, len(jwksset.Keys))
			}

			method := jwt.GetSigningMethod(tc.alg)
			token := jwt.NewWithClaims(method, jwt.StandardClaims{
				ExpiresAt: 15000,
				Issuer:    "test",
			})

			// private key is last in the array
			tokenString, err := token.SignedString(jwksset.Keys[len(jwksset.Keys)-1].Key)
			a.Nil(err)
			t.Logf("alg %s , token %s", tc.alg, tokenString)

			parser := new(jwt.Parser)
			parser.SkipClaimsValidation = true
			token, err = parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// public key is first in the array
				return jwksset.Keys[0].Key, nil
			})
			a.Nil(err)
		})
	}
}
