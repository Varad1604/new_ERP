package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/enterprise-erp/core/internal/domain"
	"github.com/enterprise-erp/core/internal/repository/postgres"
)

type FinanceUseCase struct {
	financeRepo *postgres.FinanceRepository
}

func NewFinanceUseCase(financeRepo *postgres.FinanceRepository) *FinanceUseCase {
	return &FinanceUseCase{financeRepo: financeRepo}
}

type CreateAccountRequest struct {
	Code     string             `json:"code" binding:"required"`
	Name     string             `json:"name" binding:"required"`
	Type     domain.AccountType `json:"type" binding:"required"`
	Currency string             `json:"currency" binding:"required,len=3"`
}

type LedgerLineRequest struct {
	AccountID   string  `json:"account_id" binding:"required,uuid"`
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description"`
}

type PostJournalRequest struct {
	ReferenceNumber string              `json:"reference_number" binding:"required"`
	Description     string              `json:"description" binding:"required"`
	Date            string              `json:"date" binding:"required"` // YYYY-MM-DD
	Lines           []LedgerLineRequest `json:"lines" binding:"required,min=2"`
}

func (u *FinanceUseCase) CreateAccount(ctx context.Context, req CreateAccountRequest) (*domain.Account, error) {
	acc := &domain.Account{
		Code:     req.Code,
		Name:     req.Name,
		Type:     req.Type,
		Currency: req.Currency,
	}

	err := u.financeRepo.CreateAccount(ctx, acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (u *FinanceUseCase) GetAccounts(ctx context.Context) ([]domain.Account, error) {
	return u.financeRepo.GetAccounts(ctx)
}

func (u *FinanceUseCase) PostJournalEntry(ctx context.Context, req PostJournalRequest, userID string) (*domain.JournalEntry, error) {
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	entry := &domain.JournalEntry{
		ReferenceNumber: req.ReferenceNumber,
		Description:     req.Description,
		Date:            parsedDate,
		CreatedBy:       userID,
		Lines:           make([]domain.LedgerLine, len(req.Lines)),
	}

	var sum float64
	for i, lr := range req.Lines {
		var desc *string
		if lr.Description != "" {
			desc = &lr.Description
		}

		entry.Lines[i] = domain.LedgerLine{
			AccountID:   lr.AccountID,
			Amount:      lr.Amount,
			Description: desc,
		}
		sum += lr.Amount
	}

	if sum != 0 {
		return nil, errors.New("journal entry does not balance; debits and credits must sum to zero")
	}

	err = u.financeRepo.PostJournalEntry(ctx, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (u *FinanceUseCase) GetDashboardMetrics(ctx context.Context, tenantID string) (*domain.DashboardMetrics, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}
	return u.financeRepo.GetDashboardMetrics(ctx, tenantID)
}
