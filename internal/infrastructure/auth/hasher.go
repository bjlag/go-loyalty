//go:generate mockgen -source ${GOFILE} -package mock -destination mock/hasher_mock.go

package auth

import "golang.org/x/crypto/bcrypt"

type IHasher interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashedPassword, password string) bool
}

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

func (h Hasher) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	return string(hashedPassword), err
}

func (h Hasher) ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
