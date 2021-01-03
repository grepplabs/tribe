package models

import "time"

type User struct {
	UserID            string    `db:"user_id"`
	CreatedAt         time.Time `db:"created_at,omitempty"`
	RealmID           string    `db:"realm_id"`
	Username          string    `db:"username"`
	EncryptedPassword string    `db:"encrypted_password"`
	Enabled           bool      `db:"enabled"`
	Email             string    `db:"email"`
	EmailVerified     bool      `db:"email_verified"`
}

func (User) TableName() string {
	return "tribe_user"
}

type UserList struct {
	Users []User
	Page  Page
}
