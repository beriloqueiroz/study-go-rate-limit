package usecase

import (
	"context"
	"time"
)

type FindCurrentLimiterByKeyDTO struct {
	Count    int
	UpdateAt time.Time
}

type SaveInputDTO struct {
	Count    int
	UpdateAt time.Time
}

type RateLimitRepository interface {
	FindCurrentLimiterByKey(ctx context.Context, key string) (*FindCurrentLimiterByKeyDTO, error)
	Save(ctx context.Context, input *SaveInputDTO) error
}
