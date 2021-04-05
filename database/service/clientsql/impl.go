package clientsql

import (
	"github.com/grepplabs/tribe/database/service"
	"github.com/upper/db/v4"
)

type APIImpl struct {
	kmsKeysetManager
	jwksManager
}

var _ service.API = (*APIImpl)(nil)

func NewAPIImpl(dbs db.Session) *APIImpl {
	return &APIImpl{
		kmsKeysetManager{dbs},
		jwksManager{dbs},
	}
}
