package domain

import (
	"time"
)

type AccountType string

const (
	Asset     AccountType = "ASSET"
	Liability AccountType = "LIABILITY"
	Equity    AccountType = "EQUITY"
	Revenue   AccountType = "REVENUE"
	Expense   AccountType = "EXPENSE"
)

type Account struct {
	ID        string      `json:"id" db:"id"`
	Code      string      `json:"code" db:"code"`
	Name      string      `json:"name" db:"name"`
	Type      AccountType `json:"type" db:"type"`
	Currency  string      `json:"currency" db:"currency"`
	IsActive  bool        `json:"is_active" db:"is_active"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
}

type JournalStatus string

const (
	Draft  JournalStatus = "DRAFT"
	Posted JournalStatus = "POSTED"
	Voided JournalStatus = "VOIDED"
)

type JournalEntry struct {
	ID              string        `json:"id" db:"id"`
	ReferenceNumber string        `json:"reference_number" db:"reference_number"`
	Description     string        `json:"description" db:"description"`
	Status          JournalStatus `json:"status" db:"status"`
	Date            time.Time     `json:"date" db:"date"`
	CreatedBy       string        `json:"created_by" db:"created_by"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	Lines           []LedgerLine  `json:"lines,omitempty" db:"-"`
}

type LedgerLine struct {
	ID             string    `json:"id" db:"id"`
	JournalEntryID string    `json:"journal_entry_id" db:"journal_entry_id"`
	AccountID      string    `json:"account_id" db:"account_id"`
	Amount         float64   `json:"amount" db:"amount"` // Positive = Debit, Negative = Credit
	Description    *string   `json:"description" db:"description"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
