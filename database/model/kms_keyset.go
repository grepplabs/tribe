package model

import (
	"time"
)

type KMSKeyset struct {
	KeysetID        string     `db:"id"`
	CreatedAt       *time.Time `db:"created_at,omitempty"`
	Name            string     `db:"name"`
	NextID          *string    `db:"next_id,omitempty"`
	EncryptedKeyset string     `db:"encrypted_keyset"`
	Description     *string    `db:"description,omitempty"`
}

func (KMSKeyset) TableName() string {
	return "tribe_kms_keyset"
}

type KMSKeysetList struct {
	KMSKeysets []KMSKeyset
	Page       Page
}
