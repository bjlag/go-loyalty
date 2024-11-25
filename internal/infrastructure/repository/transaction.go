//go:generate mockgen -source ${GOFILE} -package mock -destination mock/accrual_mock.go

package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/bjlag/go-loyalty/internal/model"
)

type TransactionRepo interface {
	Withdrawals(ctx context.Context, accountGUID string) ([]model.Transaction, error)
}

type TransactionPG struct {
	db *sqlx.DB
}

func NewTransactionPG(db *sqlx.DB) *TransactionPG {
	return &TransactionPG{
		db: db,
	}
}

func (r TransactionPG) Withdrawals(ctx context.Context, accountGUID string) ([]model.Transaction, error) {
	query := `
		SELECT guid, account_guid, order_number, type, sum, processed_at 
		FROM transactions 
		WHERE account_guid = $1 AND type = $2
		ORDER BY processed_at DESC
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	rows, err := stmt.QueryContext(ctx, accountGUID, model.Withdraw)
	if err != nil {
		return nil, fmt.Errorf("failed to execute a prepared query: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	var models []transaction
	for rows.Next() {
		var m transaction
		err = rows.Scan(&m.GUID, &m.AccountGUID, &m.OrderNumber, &m.Type, &m.Sum, &m.ProcessedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		models = append(models, m)
	}

	if len(models) == 0 {
		return nil, nil
	}

	result := make([]model.Transaction, 0, len(models))
	for _, m := range models {
		result = append(result, *m.export())
	}

	return result, nil
}
