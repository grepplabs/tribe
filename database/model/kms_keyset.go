package model

import (
	"time"
)

type KMSKeyset struct {
	ID              string    `db:"id" json:"id"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	EncryptedKeyset string    `db:"encrypted_keyset" json:"encrypted_keyset"`
	Description     string    `db:"description" json:"description"`
}

func (KMSKeyset) TableName() string {
	return "tribe_kms_keyset"
}

type KMSKeysetList struct {
	List []KMSKeyset `json:"list"`
	Page Page        `json:"page"`
}
