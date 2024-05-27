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
	if input.Key != "" {
		config, err := uc.configLimitRepository.FindLimitConfigByKey(ctx, input.Key)

		if err != nil {
			fmt.Println("FindLimitConfigByKey error", err)
			return nil, err
		}
		counter, err := uc.rateLimitRepository.FindCurrentLimiterByKey(ctx, input.Key)
		if err != nil {
			fmt.Println("FindCurrentLimiterByKey by key error", err)
			return nil, err
		}
		limiter := entity.NewKeyLimiter(
			entity.NewLimiterInfo(input.Key, counter.Count, config.LimitPerSecond, counter.UpdateAt, config.ExpirationTimeInMinutes, counter.StartAt),
		)

		err = uc.rateLimitRepository.Save(ctx, &SaveInputDTO{
			Count:    limiter.KeyInfo.Count,
			UpdateAt: limiter.KeyInfo.UpdateAt,
			Key:      limiter.KeyInfo.Key,
			StartAt:  limiter.KeyInfo.StartAt,
		})

		if err != nil {
			fmt.Println("Save by key error", err)
			return nil, err
		}

		return &RateLimitUseCaseOutputDto{
			Allow: !limiter.IsBlock(),
		}, nil
	}

	config, err := uc.configLimitRepository.FindLimitConfigByIp(ctx, input.Ip)
	if err != nil {
		fmt.Println("FindLimitConfigByIp error", err)
		return nil, err
	}
	counter, err := uc.rateLimitRepository.FindCurrentLimiterByKey(ctx, input.Ip)
	if err != nil {
		fmt.Println("FindCurrentLimiterByKey by ip error", err)
		return nil, err
	}
	limiter := entity.NewIpLimiter(
		entity.NewLimiterInfo(input.Ip, counter.Count, config.LimitPerSecond, counter.UpdateAt, config.ExpirationTimeInMinutes, counter.StartAt),
	)

	err = uc.rateLimitRepository.Save(ctx, &SaveInputDTO{
		Count:    limiter.IpInfo.Count,
		UpdateAt: limiter.IpInfo.UpdateAt,
		Key:      limiter.IpInfo.Key,
		StartAt:  limiter.IpInfo.StartAt,
	})

	if err != nil {
		fmt.Println("Save by ip error", err)
		return nil, err
	}

	return &RateLimitUseCaseOutputDto{
		Allow: !limiter.IsBlock(),
	}, nil

}
