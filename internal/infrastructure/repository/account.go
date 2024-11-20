//go:generate mockgen -source ${GOFILE} -package mock -destination mock/accrual_mock.go

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AccountRepo interface {
	Balance(ctx context.Context, accountGUID string) (float32, error)
}

type AccountPG struct {
	db *sqlx.DB
}

func NewAccountPG(db *sqlx.DB) *AccrualPG {
	return &AccrualPG{
		db: db,
	}
}

func (r AccrualPG) Balance(ctx context.Context, accountGUID string) (float32, error) {
	query := `SELECT balance FROM accounts WHERE guid = $1`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	var balance float32
	row := stmt.QueryRowContext(ctx, accountGUID)
	if row.Err() != nil {
		return 0, fmt.Errorf("failed to query account: %w", row.Err())
	}

	if err := row.Scan(&balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("failed to scan row: %w", err)
	}

	return balance, nil
}
