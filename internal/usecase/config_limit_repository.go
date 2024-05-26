package usecase

import (
	"context"
	"time"
)

type FindLimitConfigByKeyDTO struct {
	ExpirationTimeInMinutes time.Duration
	LimitPerSecond          int
}

type FindLimitConfigByIpDTO struct {
	ExpirationTimeInMinutes time.Duration
	LimitPerSecond          int
}

type ConfigLimitRepository interface {
	FindLimitConfigByKey(ctx context.Context, key string) (*FindLimitConfigByKeyDTO, error)
	FindLimitConfigByIp(ctx context.Context, ip string) (*FindLimitConfigByIpDTO, error)
}
