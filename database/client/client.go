package client

import (
	"github.com/grepplabs/tribe/pkg/services/resources/realms"
	"github.com/grepplabs/tribe/pkg/services/resources/users"
)

type Client interface {
	UserManager() users.Manager
	RealmManager() realms.Manager
}
