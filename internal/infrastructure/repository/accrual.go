//go:generate mockgen -source ${GOFILE} -package mock -destination mock/accrual_mock.go

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/bjlag/go-loyalty/internal/infrastructure/guid"
	"github.com/bjlag/go-loyalty/internal/model"
)

type AccrualRepo interface {
	AccrualByOrderNumber(ctx context.Context, orderNumber string) (*model.Accrual, error)
	AccrualsByUser(ctx context.Context, userGUID string) ([]model.Accrual, error)
	AccrualsInWork(ctx context.Context) ([]model.Accrual, error)
	Create(ctx context.Context, accrual *model.Accrual) error
	UpdateStatus(ctx context.Context, orderNumber string, newStatus model.AccrualStatus) error
	AddTx(ctx context.Context, accrual model.Accrual) error
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

func (r AccrualPG) AddTx(ctx context.Context, accrual model.Accrual) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// обновляем статус начисления

	query := `UPDATE accruals SET status = $1, accrual = $2 WHERE order_number = $3`
	stmtAccrual, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare update accrual query: %w", err)
	}
	defer func() {
		_ = stmtAccrual.Close()
	}()

	_, err = stmtAccrual.ExecContext(ctx, accrual.Status, accrual.Accrual, accrual.OrderNumber)
	if err != nil {
		return fmt.Errorf("failed to update accrual: %w", err)
	}

	// обновляем счет
	query = `
		INSERT INTO accounts (guid, balance, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (guid) DO UPDATE
    		SET balance    = accounts.balance + excluded.balance,
        		updated_at = excluded.updated_at;
	`
	stmtAccount, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare insert account query: %w", err)
	}
	defer func() {
		_ = stmtAccount.Close()
	}()

	_, err = stmtAccount.ExecContext(ctx, accrual.UserGUID, accrual.Accrual, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update accrual: %w", err)
	}

	// записываем транзакцию
	query = `
		INSERT INTO transactions (guid, account_guid, order_number, sum, processed_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	stmtTx, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare insert transaction query: %w", err)
	}
	defer func() {
		_ = stmtTx.Close()
	}()

	_, err = stmtTx.ExecContext(ctx, new(guid.Generator).Generate(), accrual.UserGUID, accrual.OrderNumber, accrual.Accrual, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update accrual: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
