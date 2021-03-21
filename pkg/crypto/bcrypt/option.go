package bcrypt

type Option func(*bcryptPasswordHasher) error

func WithBCryptCost(bcryptCost int) Option {
	return func(h *bcryptPasswordHasher) error {
		h.bcryptCost = bcryptCost
		return nil
	}
}
