package model

import "time"

type JWKS struct {
	ID            string    `db:"id" json:"id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	Kid           string    `db:"kid" json:"kid"`
	Alg           string    `db:"alg" json:"alg"`
	Use           string    `db:"use" json:"use"`
	KMSKeyURI     string    `db:"kms_key_uri" json:"kms_key_uri"`
	EncryptedJwks string    `db:"encrypted_jwks" json:"encrypted_jwks"`
	Description   string    `db:"description" json:"description"`
}

func (JWKS) TableName() string {
	return "tribe_jwks"
}

type JWKSList struct {
	List []JWKS `json:"list"`
	Page Page   `json:"page"`
}
