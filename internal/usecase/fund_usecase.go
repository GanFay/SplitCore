package usecase

import (
	"SplitCore/internal/domain"
	"SplitCore/internal/repository"
	"context"
)

type FundUsecase struct {
	fr repository.FundRepository
	ur repository.UserRepository
}

func NewFundUsecase(fr repository.FundRepository, ur repository.UserRepository) *FundUsecase {
	return &FundUsecase{fr: fr, ur: ur}
}

func (u *FundUsecase) GetBalance(ctx context.Context, fundID int) (*domain.Settlement, error) {
	return u.GetBalance(ctx, fundID)
}

func (u *FundUsecase) AddExpense(ctx context.Context, p domain.Purchase) error {
	return u.AddExpense(ctx, p)
}
