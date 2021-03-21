package keygen

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"github.com/pkg/errors"

	"golang.org/x/crypto/ed25519"
	"gopkg.in/square/go-jose.v2"
)

// https://tools.ietf.org/html/rfc7517
// https://tools.ietf.org/html/rfc7518

type KeygenSig interface {
	// Generate generates keypair for asymmetric algorithms specified by https://tools.ietf.org/html/rfc7518#section-3.1
	Generate(alg jose.SignatureAlgorithm) (crypto.PublicKey, crypto.PrivateKey, error)
}

func NewKeygenSig(options ...Option) KeygenSig {
	gen := &keygenSig{}
	for _, option := range options {
		option(&gen.keygen)
	}
	return gen
}

type keygenSig struct {
	keygen
}

func (r *keygenSig) verifyBits(alg jose.SignatureAlgorithm) error {
	switch alg {
	case jose.ES256, jose.ES384, jose.ES512, jose.EdDSA:
		keylen := map[jose.SignatureAlgorithm]int{
			jose.ES256: 256,
			jose.ES384: 384,
			jose.ES512: 521,
			jose.EdDSA: 256,
		}
		if r.bits != 0 && r.bits != keylen[alg] {
			return errors.New("this `alg` does not support arbitrary key length")
		}
	case jose.RS256, jose.RS384, jose.RS512, jose.PS256, jose.PS384, jose.PS512:
		if r.bits == 0 {
			r.bits = RSADefaultKeySize
		}
		if r.bits < RSAMinKeySize {
			return errors.New("too short key for RSA `alg`, 2048+ is required")
		}
	}
	return nil
}

func (r keygenSig) Generate(alg jose.SignatureAlgorithm) (crypto.PublicKey, crypto.PrivateKey, error) {
	if err := r.verifyBits(alg); err != nil {
		return nil, nil, err
	}
	switch alg {
	case jose.ES256:
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return key.Public(), key, err
	case jose.ES384:
		key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return key.Public(), key, err
	case jose.ES512:
		key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return key.Public(), key, err
	case jose.RS256, jose.RS384, jose.RS512, jose.PS256, jose.PS384, jose.PS512:
		key, err := rsa.GenerateKey(rand.Reader, r.bits)
		if err != nil {
			return nil, nil, err
		}
		return key.Public(), key, err
	case jose.EdDSA:
		pub, key, err := ed25519.GenerateKey(rand.Reader)
		return pub, key, err
	case jose.HS256, jose.HS384, jose.HS512:
		return nil, nil, errors.Errorf("symmetric alg=%s used to generate keypair", alg)
	default:
		return nil, nil, errors.Errorf("unknown alg=%s for use=sig", alg)
	}
}
