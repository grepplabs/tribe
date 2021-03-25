package model

import "time"

type KMSKeyset struct {
	Id              string    `db:"id"`
	CreatedAt       time.Time `db:"created_at,omitempty"`
	Name            string    `db:"name"`
	NextId          string    `db:"next_id,omitempty"`
	EncryptedKeyset string    `db:"encrypted_keyset"`
	Description     string    `db:"description"`
}

func (KMSKeyset) TableName() string {
	return "tribe_kms_keyset"
}

type KMSKeysetList struct {
	KMSKeysets []KMSKeyset
}
