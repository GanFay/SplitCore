package usecase

import (
	"SplitCore/internal/domain"
	"SplitCore/internal/repository"
	"context"
)

type StatesUsecase struct {
	redisRep repository.RedisRepository
}

func NewStateUsecase(redisRep repository.RedisRepository) *StatesUsecase {
	return &StatesUsecase{redisRep: redisRep}
}

func (r *StatesUsecase) GetUserCtx(ctx context.Context, userID int64) (*domain.UserContext, error) {
	return r.redisRep.GetUserCtx(ctx, userID)
}

func (r *StatesUsecase) SaveUserCtx(ctx context.Context, userID int64, value *domain.UserContext) error {
	return r.redisRep.SaveUserCtx(ctx, userID, value)
}
