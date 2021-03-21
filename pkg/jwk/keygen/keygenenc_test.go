package keygen

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/square/go-jose.v2"
	"testing"
)

func TestKeygenEnc(t *testing.T) {
	// only asymmetric key from https://tools.ietf.org/html/rfc7518#section-3.1
	tests := []struct {
		alg  jose.KeyAlgorithm
		bits int
	}{
		{alg: jose.RSA1_5},
		{alg: jose.RSA1_5, bits: 2048},
		{alg: jose.RSA1_5, bits: 3072},
		{alg: jose.RSA_OAEP},
		{alg: jose.RSA_OAEP, bits: 2048},
		{alg: jose.RSA_OAEP, bits: 3072},
		{alg: jose.RSA_OAEP_256},
		{alg: jose.RSA_OAEP_256}, // RSA_OAEP_256 + A128CBC_HS256
		{alg: jose.RSA_OAEP_256, bits: 2048},
		{alg: jose.RSA_OAEP_256, bits: 3072},

		{alg: jose.ECDH_ES},
		{alg: jose.ECDH_ES, bits: 256},
		{alg: jose.ECDH_ES, bits: 384},
		{alg: jose.ECDH_ES, bits: 521},

		{alg: jose.ECDH_ES_A128KW},
		{alg: jose.ECDH_ES_A128KW, bits: 256},
		{alg: jose.ECDH_ES_A128KW, bits: 384},
		{alg: jose.ECDH_ES_A128KW, bits: 521},

		{alg: jose.ECDH_ES_A192KW},
		{alg: jose.ECDH_ES_A192KW, bits: 256},
		{alg: jose.ECDH_ES_A192KW, bits: 384},
		{alg: jose.ECDH_ES_A192KW, bits: 521},

		{alg: jose.ECDH_ES_A256KW},
		{alg: jose.ECDH_ES_A256KW, bits: 256},
		{alg: jose.ECDH_ES_A256KW, bits: 384},
		{alg: jose.ECDH_ES_A256KW, bits: 521},
	}
	for _, tc := range tests {
		name := string(tc.alg)
		if tc.bits != 0 {
			name = fmt.Sprintf("%s-%d", tc.alg, tc.bits)
		}
		t.Run(name, func(t *testing.T) {
			a := assert.New(t)
			gen := NewKeygenEnc(WithBits(tc.bits))
			publicKey, privateKey, err := gen.Generate(tc.alg)
			a.Nil(err)
			a.NotNil(publicKey)
			a.NotNil(privateKey)

			enc := jose.A128CBC_HS256
			crypter, err := jose.NewEncrypter(enc, jose.Recipient{Algorithm: tc.alg, Key: publicKey}, nil)
			a.Nil(err)

			input := []byte("Lorem ipsum dolor sit amet")
			obj, err := crypter.Encrypt(input)
			a.Nil(err)

			msg := obj.FullSerialize()
			_ = msg

			obj, err = jose.ParseEncrypted(msg)
			a.Nil(err)
			plaintext, err := obj.Decrypt(privateKey)
			a.Equal(input, plaintext)
		})
	}
}
