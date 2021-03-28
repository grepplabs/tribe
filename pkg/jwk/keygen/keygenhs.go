package keygen

import (
	"crypto/rand"
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
	"io"
)

type KeygenHs interface {
	Generate(alg jose.SignatureAlgorithm) ([]byte, error)
}

func NewKeygenHs() KeygenHs {
	return &keygenHs{}
}

type keygenHs struct {
	key string
}

func (r keygenHs) Generate(alg jose.SignatureAlgorithm) ([]byte, error) {
	var key []byte
	switch alg {
	case jose.HS256:
		key = make([]byte, 32)
	case jose.HS384:
		key = make([]byte, 48)
	case jose.HS512:
		key = make([]byte, 64)
	default:
		return nil, errors.Errorf("unknown symmetric `alg` %s", alg)
	}
	_, err := io.ReadFull(rand.Reader, key[:])
	return key, err
}
