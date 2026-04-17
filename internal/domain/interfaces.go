package domain

import "context"

type FundUsecase interface {
	GetBalance(ctx context.Context, fundID int) (*Settlement, error)
	AddExpense(ctx context.Context, p Purchase) error
}
