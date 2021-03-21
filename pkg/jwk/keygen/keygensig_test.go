package keygen

import (
	"fmt"
	"testing"

	"github.com/form3tech-oss/jwt-go"
	"github.com/stretchr/testify/assert"
	"gopkg.in/square/go-jose.v2"
)

func TestKeygenSig(t *testing.T) {
	// only asymmetric key from https://tools.ietf.org/html/rfc7518#section-3.1
	tests := []struct {
		alg  jose.SignatureAlgorithm
		bits int
	}{
		{alg: jose.RS256},
		{alg: jose.RS384},
		{alg: jose.RS512},
		{alg: jose.ES256},
		{alg: jose.ES384},
		{alg: jose.ES512},
		{alg: jose.PS256},
		{alg: jose.PS384},
		{alg: jose.PS512},
		{alg: jose.RS256, bits: 2048},
		{alg: jose.RS384, bits: 2048},
		{alg: jose.RS512, bits: 2048},
		{alg: jose.RS256, bits: 3072},
		{alg: jose.RS384, bits: 3072},
		{alg: jose.RS512, bits: 3072},
	}
	for _, tc := range tests {
		name := string(tc.alg)
		if tc.bits != 0 {
			name = fmt.Sprintf("%s-%d", tc.alg, tc.bits)
		}
		t.Run(name, func(t *testing.T) {
			a := assert.New(t)
			gen := NewKeygenSig(WithBits(tc.bits))
			publicKey, privateKey, err := gen.Generate(tc.alg)
			a.Nil(err)
			a.NotNil(publicKey)
			a.NotNil(privateKey)

			// sign and verify with generated keys
			method := jwt.GetSigningMethod(string(tc.alg))
			token := jwt.NewWithClaims(method, jwt.StandardClaims{
				ExpiresAt: 15000,
				Issuer:    "test",
			})

			ss, err := token.SignedString(privateKey)
			a.Nil(err)
			parser := new(jwt.Parser)
			parser.SkipClaimsValidation = true
			token, err = parser.Parse(ss, func(token *jwt.Token) (interface{}, error) {
				return publicKey, nil
			})
			a.Nil(err)

			signer, err := jose.NewSigner(jose.SigningKey{Algorithm: tc.alg, Key: privateKey}, &jose.SignerOptions{})
			a.Nil(err)
			input := []byte("Lorem ipsum dolor sit amet")
			_, err = signer.Sign(input)
			a.Nil(err)
		})
	}
}
