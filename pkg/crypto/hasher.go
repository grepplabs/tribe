package crypto

type PasswordHasher interface {

	// HashPassword creates a hash from plain password or returns an error.
	HashPassword(password string) (string, error)

	// VerifyHashedPassword compares plain password with a hash and returns an error  if the two do not match.
	VerifyHashedPassword(password, hash string) bool
}
