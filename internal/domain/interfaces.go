package domain

import (
	"context"
)

type FundUsecase interface {
	GetBalance(ctx context.Context, fundID int) (*Settlement, error)
	AddExpense(ctx context.Context, fundID int, id int64, text string) (*Purchase, error)

	CreateFund(ctx context.Context, fund *Fund) (*Fund, error)
	GetInfo(ctx context.Context, reqFund *Fund) (*Fund, error)
	GetByUserID(ctx context.Context, userID int64, limit int, offset int) ([]Fund, error)
	AddMember(ctx context.Context, fund *Fund, userID int64) error
	IsMember(ctx context.Context, fundID int, userID int64) (bool, error)
	GetMembers(ctx context.Context, fundID int) ([]User, error)

	GetPurchasesByFundPagination(ctx context.Context, fundID int, limit int, offset int) ([]Purchase, error)
	CreatePurchase(ctx context.Context, purchase *Purchase) error
}

type UserUsecase interface {
	CreateRealUser(ctx context.Context, tgID int64, username string, firstName string) (int64, error)
	GetUser(ctx context.Context, id int64) (*User, error)
	CreateVirtualUser(ctx context.Context, firstName string) (int64, error)
}

type StatesUsecase interface {
	GetUserCtx(ctx context.Context, userID int64) (*UserContext, error)
	SaveUserCtx(ctx context.Context, userID int64, value *UserContext) error
}
