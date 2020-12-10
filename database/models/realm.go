package models

import "time"

type Realm struct {
	RealmID     string    `db:"realm_id"`
	CreatedAt   time.Time `db:"created_at"`
	Description string    `db:"description"`
}

func (Realm) TableName() string {
	return "tribe_realm"
}
