package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/bjlag/go-loyalty/internal/model"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
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

func (r UserRepository) Insert(ctx context.Context, user *model.User) error {
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
