package repository

import (
	"context"
	"time"

	"github.com/beriloqueiroz/study-go-rate-limit/internal/usecase"
)

type RateLimitRepositoryImpl struct {
}

func (rr *RateLimitRepositoryImpl) FindCurrentLimiterByKey(ctx context.Context, key string) (*usecase.FindCurrentLimiterByKeyDTO, error) {
	return &usecase.FindCurrentLimiterByKeyDTO{
		Count:    0,
		UpdateAt: time.Now(),
	}, nil
}

func (rr *RateLimitRepositoryImpl) Save(ctx context.Context, input *usecase.SaveInputDTO) error {
	return nil
}
