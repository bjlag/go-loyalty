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
	FindByLogin(ctx context.Context, login string) (*model.User, error)
	Insert(ctx context.Context, user *model.User) error
}

type UserPG struct {
	db *sqlx.DB
}

func NewUserPG(db *sqlx.DB) *UserPG {
	return &UserPG{
		db: db,
	}
}

func (r UserPG) FindByLogin(ctx context.Context, login string) (*model.User, error) {
	query := "SELECT guid, login, password FROM users WHERE login = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	var m user
	row := stmt.QueryRowContext(ctx, login)
	err = row.Scan(&m.GUID, &m.Login, &m.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return m.export(), nil
}

func (r UserPG) Insert(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (guid, login, password) VALUES ($1, $2, $3)`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	//m := userFromModel(user)
	_, err = stmt.ExecContext(ctx, user.GUID, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}
