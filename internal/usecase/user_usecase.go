package usecase

import (
	"context"

	"github.com/ganfay/split-core/internal/domain"
	"github.com/ganfay/split-core/internal/repository"
)

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) domain.UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) CreateRealUser(ctx context.Context, tgID int64, username string, firstName string) (int64, error) {
	return u.repo.CreateRealUser(ctx, tgID, username, firstName)
}

func (u *userUsecase) GetUser(ctx context.Context, tgID int64) (*domain.User, error) {
	return u.repo.GetUser(ctx, tgID)
}

func (u *userUsecase) CreateVirtualUser(ctx context.Context, firstName string) (int64, error) {
	return u.repo.CreateVirtualUser(ctx, firstName)
}
