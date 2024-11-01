package auth

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
	cost int
}

type Option func(*Hasher)

func WithCost(cost int) Option {
	return func(h *Hasher) {
		h.cost = cost
	}
}

func NewHasher(opts ...Option) *Hasher {
	h := &Hasher{
		cost: bcrypt.MinCost,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(hashedPassword), err
}
