//go:generate mockgen -source ${GOFILE} -package mock -destination mock/user_mock.go

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/bjlag/go-loyalty/internal/model"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Insert(ctx context.Context, user *model.User) error
}

type UserPG struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserPG {
	return &UserPG{
		db: db,
	}
}

func (r UserPG) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	var m user
	row := stmt.QueryRowContext(ctx, email)
	err = row.Scan(&m.GUID, &m.Email, &m.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return m.export(), nil
}

func (r UserPG) Insert(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (guid, email, password) VALUES ($1, $2, $3)`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	//m := userFromModel(user)
	_, err = stmt.ExecContext(ctx, user.GUID, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}
