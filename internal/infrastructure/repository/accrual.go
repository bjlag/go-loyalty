package repository

import (
	"context"
	"fmt"
	"github.com/bjlag/go-loyalty/internal/model"
	"github.com/jmoiron/sqlx"
)

type AccrualRepository interface {
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

func (r AccrualPG) Insert(ctx context.Context, accrual *model.Accrual) error {
	query := `INSERT INTO accruals (number, user_guid, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4, $5)`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.ExecContext(ctx, accrual.Number, accrual.UserGUID, accrual.Status, accrual.Accrual, accrual.UploadedAt)
	if err != nil {
		return fmt.Errorf("failed to save accrual: %w", err)
	}

	return nil
}
