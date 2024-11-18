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

type AccrualRepo interface {
	AccrualByOrderNumber(ctx context.Context, orderNumber string) (*model.Accrual, error)
	AccrualsByUser(ctx context.Context, userGUID string) ([]model.Accrual, error)
	AccrualsInWork(ctx context.Context) ([]model.Accrual, error)

	Create(ctx context.Context, accrual *model.Accrual) error
	UpdateStatus(ctx context.Context, orderNumber string, newStatus model.AccrualStatus) error
	Add(ctx context.Context, accrual model.Accrual, account model.Account, transaction model.Transaction) error
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

func (r AccrualPG) AccrualsByUser(ctx context.Context, orderNumber string) ([]model.Accrual, error) {
	query := "SELECT order_number, user_guid, status, accrual, uploaded_at FROM accruals WHERE user_guid = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	rows, err := stmt.QueryContext(ctx, orderNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to execute a prepared query: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	var models []accrual
	for rows.Next() {
		var m accrual
		err = rows.Scan(&m.OrderNumber, &m.UserGUID, &m.Status, &m.Accrual, &m.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		models = append(models, m)
	}

	if len(models) == 0 {
		return nil, nil
	}

	result := make([]model.Accrual, 0, len(models))
	for _, m := range models {
		result = append(result, *m.export())
	}

	return result, nil
}

func (r AccrualPG) AccrualsInWork(ctx context.Context) ([]model.Accrual, error) {
	query := "SELECT order_number, user_guid, status, accrual, uploaded_at FROM accruals WHERE status = $1 OR status = $2"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	rows, err := stmt.QueryContext(ctx, model.New, model.Processing)
	if err != nil {
		return nil, fmt.Errorf("failed to execute a prepared query: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	var models []accrual
	for rows.Next() {
		var m accrual
		err = rows.Scan(&m.OrderNumber, &m.UserGUID, &m.Status, &m.Accrual, &m.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		models = append(models, m)
	}

	if len(models) == 0 {
		return nil, nil
	}

	result := make([]model.Accrual, 0, len(models))
	for _, m := range models {
		result = append(result, *m.export())
	}

	return result, nil
}

func (r AccrualPG) Create(ctx context.Context, accrual *model.Accrual) error {
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

func (r AccrualPG) UpdateStatus(ctx context.Context, orderNumber string, newStatus model.AccrualStatus) error {
	query := `UPDATE accruals SET status = $1 WHERE order_number = $2`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.ExecContext(ctx, newStatus, orderNumber)
	if err != nil {
		return fmt.Errorf("failed to update accrual: %w", err)
	}

	return nil
}

func (r AccrualPG) Add(ctx context.Context, accrual model.Accrual, account model.Account, transaction model.Transaction) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	err = updateAccrualTx(tx, accrual)
	if err != nil {
		return err
	}

	err = updateAccountTx(tx, account)
	if err != nil {
		return err
	}

	err = addTransaction(tx, transaction)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func updateAccrualTx(tx *sql.Tx, model model.Accrual) error {
	query := `UPDATE accruals SET status = $1, accrual = $2 WHERE order_number = $3`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare update accrual query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(model.Status, model.Accrual, model.OrderNumber)
	if err != nil {
		return fmt.Errorf("failed to update accrual: %w", err)
	}

	return nil
}

func updateAccountTx(tx *sql.Tx, model model.Account) error {
	query := `
		INSERT INTO accounts (guid, balance, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (guid) DO UPDATE
    		SET balance    = accounts.balance + excluded.balance,
        		updated_at = excluded.updated_at;
	`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare insert account query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(model.GUID, model.Balance, model.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	return nil
}

func addTransaction(tx *sql.Tx, model model.Transaction) error {
	query := `
		INSERT INTO transactions (guid, account_guid, order_number, sum, processed_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare insert transaction query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(model.GUID, model.AccountGUID, model.OrderNumber, model.Sum, model.ProcessedAt)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	return nil
}
