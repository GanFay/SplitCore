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

func (u *userUsecase) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return u.repo.CreateRealUser(ctx, user)
}

func (u *userUsecase) GetUser(ctx context.Context, tgID int64) (*domain.User, error) {
	return u.repo.GetUser(ctx, tgID)
}
