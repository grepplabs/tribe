package model

import "time"

type OidcJWKS struct {
	ID             string    `db:"id" json:"id"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	CurrentJwksID  string    `db:"current_jwks_id" json:"current_jwks_id"`
	NextJwksID     string    `db:"next_jwks_id" json:"next_jwks_id"`
	PreviousJwksID *string   `db:"previous_jwks_id" json:"previous_jwks_id,omitempty"`
	RotationMode   int       `db:"rotation_mode" json:"rotation_mode"`
	RotationPeriod int       `db:"rotation_period" json:"rotation_period"`
	LastRotated    time.Time `db:"last_rotated" json:"last_rotated"`
	Description    string    `db:"description" json:"description"`
	Version        int       `db:"version" json:"version"`
}

func (OidcJWKS) TableName() string {
	return "tribe_oidc_jwks"
}

type OidcJWKSList struct {
	List []OidcJWKS `json:"list"`
	Page Page       `json:"page"`
}
