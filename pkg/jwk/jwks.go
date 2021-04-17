package jwk

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/grepplabs/tribe/pkg/jwk/keygen"
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
)

const (
	privateKeyIDPrefix = "private"
	publicKeyIDPrefix  = "public"
)

type JWKSGenerator interface {
	Generate(id, alg, use string) (*jose.JSONWebKeySet, error)
}

func NewJWKSGenerator() JWKSGenerator {
	return &jwksGenerator{}
}

type jwksGenerator struct {
}

func (g jwksGenerator) Generate(id, alg, use string) (*jose.JSONWebKeySet, error) {
	// https://tools.ietf.org/html/rfc7518#page-6
	if use != "enc" && use != "sig" {
		return nil, errors.Errorf("unsupported intend of use %s", use)
	}
	if id == "" {
		id = uuid.NewString()
	}
	switch alg {
	case "HS256", "HS384", "HS512":
		key, err := keygen.NewKeygenHs().Generate(jose.SignatureAlgorithm(alg))
		if err != nil {
			return nil, err
		}
		return &jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{
				{
					Algorithm: alg,
					Use:       use,
					Key:       key,
					KeyID:     id,
				},
			},
		}, nil
	case "RS256", "RS384", "RS512", "ES256", "ES384", "ES512", "PS256", "PS384", "PS512":
		//TODO: generate with self signed certs
		publicKey, privateKey, err := keygen.NewKeygenSig().Generate(jose.SignatureAlgorithm(alg))
		if err != nil {
			return nil, err
		}
		return &jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{
				{
					Algorithm: alg,
					Use:       use,
					Key:       publicKey,
					KeyID:     g.keyID(publicKeyIDPrefix, id),
				},
				{
					Algorithm: alg,
					Use:       use,
					Key:       privateKey,
					KeyID:     g.keyID(privateKeyIDPrefix, id),
				},
			},
		}, nil
	case "none":
		return nil, errors.New("unsecure 'none' algorithm is not supported")
	default:
		return nil, errors.Errorf("unsupported alg: %s", alg)
	}
}

func (g jwksGenerator) keyID(prefix, id string) string {
	return fmt.Sprintf("%s-%s", prefix, id)
}
