package usecase

import (
	"SplitCore/internal/domain"
	"SplitCore/internal/repository"
	"context"
)

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) domain.UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return u.repo.CreateUser(ctx, user)
}

func (u *userUsecase) GetUser(ctx context.Context, tgID int64) (*domain.User, error) {
	return u.repo.GetUser(ctx, tgID)
}
