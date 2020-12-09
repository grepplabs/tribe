package models

type Realm struct {
	RealmID     string `db:"realm_id"`
	Description string `db:"description"`
}

func (Realm) TableName() string {
	return "tribe_realm"
}
