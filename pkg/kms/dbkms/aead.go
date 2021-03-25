package dbkms

import (
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/tink"
)

type GetKeysetFunc func() (*keyset.Handle, error)

type kmsAEAD struct {
	getKeysetFunc GetKeysetFunc
}

var _ tink.AEAD = (*kmsAEAD)(nil)

func NewAEAD(getKeyset GetKeysetFunc) tink.AEAD {
	return &kmsAEAD{getKeysetFunc: getKeyset}
}

// Encrypt AEAD encrypts the plaintext data and uses addtionaldata from authentication.
func (r *kmsAEAD) Encrypt(plaintext, additionalData []byte) ([]byte, error) {
	h, err := r.getKeysetFunc()
	if err != nil {
		return nil, err
	}
	a, err := aead.New(h)
	if err != nil {
		return nil, err
	}
	return a.Encrypt(plaintext, additionalData)
}

// Decrypt AEAD decrypts the data and verified the additional data.
func (r *kmsAEAD) Decrypt(ciphertext, additionalData []byte) ([]byte, error) {
	h, err := r.getKeysetFunc()
	if err != nil {
		return nil, err
	}
	a, err := aead.New(h)
	if err != nil {
		return nil, err
	}
	return a.Decrypt(ciphertext, additionalData)
}
