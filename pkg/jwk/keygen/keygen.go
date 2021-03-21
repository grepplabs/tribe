package keygen

// https://tools.ietf.org/html/rfc7517
// https://tools.ietf.org/html/rfc7518

const (
	RSADefaultKeySize = 4096
	RSAMinKeySize     = 2048
)

type keygen struct {
	bits int
}

type Option func(*keygen)

func WithBits(bits int) Option {
	return func(h *keygen) {
		h.bits = bits
	}
}
