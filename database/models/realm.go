package models

import (
	"time"
)

type Realm struct {
	RealmID     string    `db:"realm_id"`
	CreatedAt   time.Time `db:"created_at,omitempty"`
	Description string    `db:"description"`
}

type RealmList struct {
	Realms []Realm
	Page   Page
}

type Page struct {
	Offset *int64
	Limit  *int64
	Total  uint64
}

func (Realm) TableName() string {
	return "tribe_realm"
}
