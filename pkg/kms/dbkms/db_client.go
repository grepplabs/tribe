package dbkms

import (
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/tink"
)

type dbClient struct {
}

func (c *dbClient) GetAEAD(keyname string) (tink.AEAD, error) {
	return NewAEAD(func() (*keyset.Handle, error) {
		//TODO: implement me
		return nil, nil
	}), nil
}
