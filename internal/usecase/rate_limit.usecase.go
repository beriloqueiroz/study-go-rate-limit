package usecase

import (
	"context"
	"fmt"

	"github.com/beriloqueiroz/study-go-rate-limit/internal/entity"
)

type RateLimitUseCase struct {
	rateLimitRepository   RateLimitRepository
	configLimitRepository ConfigLimitRepository
}

func NewRateLimitUseCase(rateLimitRepository RateLimitRepository, configLimitRepository ConfigLimitRepository) *RateLimitUseCase {
	return &RateLimitUseCase{
		rateLimitRepository:   rateLimitRepository,
		configLimitRepository: configLimitRepository,
	}
}

type RateLimitUseCaseInputDto struct {
	Ip  string
	Key string
}

type RateLimitUseCaseOutputDto struct {
	Allow bool
}

func (uc *RateLimitUseCase) Execute(ctx context.Context, input RateLimitUseCaseInputDto) (*RateLimitUseCaseOutputDto, error) {
	fmt.Println(input)
	if input.Key != "" {
		config, err := uc.configLimitRepository.FindLimitConfigByKey(ctx, input.Key)
		fmt.Println(config)

		if err != nil {
			return nil, err
		}
		counter, err := uc.rateLimitRepository.FindCurrentLimiterByKey(ctx, input.Key)
		if err != nil {
			return nil, err
		}
		limiter := entity.NewKeyLimiter(
			*entity.NewLimiterInfo(input.Key, counter.Count, config.LimitPerSecond, counter.UpdateAt, config.ExpirationTimeInMinutes),
		)
		return &RateLimitUseCaseOutputDto{
			Allow: !limiter.IsBlock(),
		}, nil
	}

	config, err := uc.configLimitRepository.FindLimitConfigByIp(ctx, input.Ip)
	if err != nil {
		return nil, err
	}
	counter, err := uc.rateLimitRepository.FindCurrentLimiterByKey(ctx, input.Ip)
	if err != nil {
		return nil, err
	}
	limiter := entity.NewIpLimiter(
		*entity.NewLimiterInfo(input.Ip, counter.Count, config.LimitPerSecond, counter.UpdateAt, config.ExpirationTimeInMinutes),
	)
	return &RateLimitUseCaseOutputDto{
		Allow: limiter.IsBlock(),
	}, nil

}
