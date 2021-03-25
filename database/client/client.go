package client

import "github.com/grepplabs/tribe/database/service"

type Client interface {
	API() service.API
}
