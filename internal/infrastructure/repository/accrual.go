//go:generate mockgen -source ${GOFILE} -package mock -destination mock/accrual_mock.go

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/bjlag/go-loyalty/internal/model"
)

type AccrualRepository interface {
	AccrualByOrderNumber(ctx context.Context, orderNumber string) (*model.Accrual, error)
	Insert(ctx context.Context, accrual *model.Accrual) error
}

type AccrualPG struct {
	db *sqlx.DB
}

func NewAccrualPG(db *sqlx.DB) *AccrualPG {
	return &AccrualPG{
		db: db,
	}
}

func (r AccrualPG) AccrualByOrderNumber(ctx context.Context, orderNumber string) (*model.Accrual, error) {
	query := "SELECT order_number, user_guid, status, accrual, uploaded_at FROM accruals WHERE order_number = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	var m accrual
	row := stmt.QueryRowContext(ctx, orderNumber)
	err = row.Scan(&m.OrderNumber, &m.UserGUID, &m.Status, &m.Accrual, &m.UploadedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return m.export(), nil
}

func (r AccrualPG) Insert(ctx context.Context, accrual *model.Accrual) error {
	query := `INSERT INTO accruals (order_number, user_guid, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4, $5)`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.ExecContext(ctx, accrual.OrderNumber, accrual.UserGUID, accrual.Status, accrual.Accrual, accrual.UploadedAt)
	if err != nil {
		return fmt.Errorf("failed to save accrual: %w", err)
	}

	return nil
}
