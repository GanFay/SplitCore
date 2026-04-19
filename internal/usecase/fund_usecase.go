package usecase

import (
	"SplitCore/internal/domain"
	"SplitCore/internal/pkg/utils"
	"SplitCore/internal/repository"
	"context"
	"errors"

	tele "gopkg.in/telebot.v4"
)

type FundUsecase struct {
	fundRepository     repository.FundRepository
	purchaseRepository repository.PurchaseRepository
	userRepository     repository.UserRepository
}

func NewFundUsecase(fr repository.FundRepository, pr repository.PurchaseRepository) *FundUsecase {
	return &FundUsecase{fundRepository: fr, purchaseRepository: pr}
}

func (u *FundUsecase) GetBalance(ctx context.Context, fundID int) (*domain.Settlement, error) {
	return u.GetBalance(ctx, fundID)
}

func (u *FundUsecase) AddExpense(ctx context.Context, ctxInfoAboutPurchase tele.Context, fundID int) (*domain.Purchase, error) {
	isMember, err := u.fundRepository.IsMember(ctx, fundID, ctxInfoAboutPurchase.Sender().ID)
	if err != nil || !isMember {
		return nil, err
	}
	cost, desc, err := utils.ParsePurchase(ctxInfoAboutPurchase.Text())
	if err != nil {
		return nil, err
	}
	if cost <= 0 {
		return nil, errors.New("invalid amount")
	}

	purchase := &domain.Purchase{
		FundID:      fundID,
		PayerID:     ctxInfoAboutPurchase.Sender().ID,
		Amount:      cost,
		Description: desc,
	}
	err = u.purchaseRepository.CreatePurchase(ctx, purchase)
	if err != nil {
		return nil, err
	}
	return purchase, nil
}

func (u *FundUsecase) CreateFund(ctx context.Context, fund *domain.Fund) (*domain.Fund, error) {
	return u.fundRepository.CreateFund(ctx, fund)
}

func (u *FundUsecase) GetInfo(ctx context.Context, reqFund *domain.Fund) (*domain.Fund, error) {
	return u.fundRepository.GetInfo(ctx, reqFund)
}

func (u *FundUsecase) GetByUserID(ctx context.Context, userID int64, limit int, offset int) ([]domain.Fund, error) {
	return u.fundRepository.GetByUserID(ctx, userID, limit, offset)
}

func (u *FundUsecase) AddMember(ctx context.Context, fund *domain.Fund, userID int64) error {
	return u.fundRepository.AddMember(ctx, fund, userID)
}

func (u *FundUsecase) CreateMember(ctx context.Context, fund *domain.Fund, userID int64) error {
	return u.fundRepository.AddMember(ctx, fund, userID)
}

func (u *FundUsecase) IsMember(ctx context.Context, fundID int, userID int64) (bool, error) {
	return u.fundRepository.IsMember(ctx, fundID, userID)
}

func (u *FundUsecase) GetPurchasesByFund(ctx context.Context, fund *domain.Fund) ([]domain.Purchase, error) {
	return u.purchaseRepository.GetPurchasesByFund(ctx, fund)
}

func (u *FundUsecase) CreatePurchase(ctx context.Context, purchase *domain.Purchase) error {
	return u.purchaseRepository.CreatePurchase(ctx, purchase)
}
