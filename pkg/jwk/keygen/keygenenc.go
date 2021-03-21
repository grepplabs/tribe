package keygen

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
)

type KeygenEnc interface {
	Generate(alg jose.KeyAlgorithm) (crypto.PublicKey, crypto.PrivateKey, error)
}

func NewKeygenEnc(options ...Option) KeygenEnc {
	gen := &keygenEnc{}
	for _, option := range options {
		option(&gen.keygen)
	}
	return gen
}

type keygenEnc struct {
	keygen
}

func (r *keygenEnc) verifyBits(alg jose.KeyAlgorithm) error {
	switch alg {
	case jose.RSA1_5, jose.RSA_OAEP, jose.RSA_OAEP_256:
		if r.bits == 0 {
			r.bits = RSADefaultKeySize
		}
		if r.bits < RSAMinKeySize {
			return errors.New("too short key for RSA `alg`, 2048+ is required")
		}
	}
	return nil
}
func (r keygenEnc) Generate(alg jose.KeyAlgorithm) (crypto.PublicKey, crypto.PrivateKey, error) {
	if err := r.verifyBits(alg); err != nil {
		return nil, nil, err
	}
	switch alg {
	case jose.RSA1_5, jose.RSA_OAEP, jose.RSA_OAEP_256:
		key, err := rsa.GenerateKey(rand.Reader, r.bits)
		if err != nil {
			return nil, nil, err
		}
		return key.Public(), key, err
	case jose.ECDH_ES, jose.ECDH_ES_A128KW, jose.ECDH_ES_A192KW, jose.ECDH_ES_A256KW:
		var crv elliptic.Curve
		switch r.bits {
		case 0, 256:
			crv = elliptic.P256()
		case 384:
			crv = elliptic.P384()
		case 521:
			crv = elliptic.P521()
		default:
			return nil, nil, errors.New("unknown elliptic curve bit length, use one of 256, 384, 521")
		}
		key, err := ecdsa.GenerateKey(crv, rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return key.Public(), key, err
	default:
		return nil, nil, errors.New("unknown `alg` for `use` = `enc`")
	}
}
