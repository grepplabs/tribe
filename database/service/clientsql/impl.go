package clientsql

import (
	"github.com/grepplabs/tribe/database/service"
	"github.com/upper/db/v4"
)

type APIImpl struct {
	KMSKeysetManager
	JWKSManager
}

var _ service.API = (*APIImpl)(nil)

func NewAPIImpl(dbs db.Session) *APIImpl {
	return &APIImpl{
		KMSKeysetManager{dbs},
		JWKSManager{dbs},
	}
}
