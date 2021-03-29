package model

import "time"

type JWKS struct {
	ID            string     `db:"id"`
	CreatedAt     *time.Time `db:"created_at,omitempty"`
	Kid           string     `db:"kid"`
	Alg           string     `db:"alg"`
	Use           string     `db:"use"`
	KMSKeysetURI  string     `db:"kms_keyset_uri"`
	EncryptedJwks string     `db:"encrypted_jwks"`
	Description   *string    `db:"description,omitempty"`
}

func (JWKS) TableName() string {
	return "tribe_jwks"
}

type JWKSList struct {
	List []JWKS
	Page Page
}
