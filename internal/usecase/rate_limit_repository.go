package usecase

import (
	"context"
	"time"
)

type FindCurrentLimiterByKeyDTO struct {
	Count    int
	UpdateAt time.Time
	StartAt  time.Time
}

type SaveInputDTO struct {
	Key      string
	Count    int
	UpdateAt time.Time
	StartAt  time.Time
}

type RateLimitRepository interface {
	FindCurrentLimiterByKey(ctx context.Context, key string) (*FindCurrentLimiterByKeyDTO, error)
	Save(ctx context.Context, input *SaveInputDTO) error
}
