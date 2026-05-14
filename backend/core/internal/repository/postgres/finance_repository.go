package postgres

import (
	"context"
	"errors"

	"github.com/enterprise-erp/core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type FinanceRepository struct {
	db *sqlx.DB
}

func NewFinanceRepository(db *sqlx.DB) *FinanceRepository {
	return &FinanceRepository{db: db}
}

func (r *FinanceRepository) CreateAccount(ctx context.Context, acc *domain.Account) error {
	query := `
		INSERT INTO accounts (code, name, type, currency)
		VALUES ($1, $2, $3, $4)
		RETURNING id, is_active, created_at
	`
	return r.db.QueryRowxContext(ctx, query, acc.Code, acc.Name, acc.Type, acc.Currency).
		Scan(&acc.ID, &acc.IsActive, &acc.CreatedAt)
}

func (r *FinanceRepository) GetAccounts(ctx context.Context) ([]domain.Account, error) {
	var accounts []domain.Account
	query := `SELECT * FROM accounts ORDER BY code ASC`
	err := r.db.SelectContext(ctx, &accounts, query)
	return accounts, err
}

func (r *FinanceRepository) PostJournalEntry(ctx context.Context, entry *domain.JournalEntry) error {
	// Start a database transaction for atomic double-entry bookkeeping
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. Insert Header
	queryHeader := `
		INSERT INTO journal_entries (reference_number, description, status, date, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err = tx.QueryRowxContext(ctx, queryHeader,
		entry.ReferenceNumber, entry.Description, domain.Posted, entry.Date, entry.CreatedBy).
		Scan(&entry.ID, &entry.CreatedAt)
	if err != nil {
		return err
	}
	entry.Status = domain.Posted

	// 2. Insert Lines & Validate Balance
	var totalBalance float64
	queryLine := `
		INSERT INTO ledger_lines (journal_entry_id, account_id, amount, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	for i := range entry.Lines {
		line := &entry.Lines[i]
		line.JournalEntryID = entry.ID
		totalBalance += line.Amount

		err = tx.QueryRowxContext(ctx, queryLine,
			line.JournalEntryID, line.AccountID, line.Amount, line.Description).
			Scan(&line.ID, &line.CreatedAt)
		if err != nil {
			return err
		}
	}

	// 3. Absolute Financial Integrity Check
	// Accounting Equation: Sum of Debits and Credits must exactly equal 0
	if totalBalance != 0.0 {
		err = errors.New("journal entry does not balance; transaction rejected")
		return err
	}

	// Commit Transaction
	err = tx.Commit()
	return err
}
