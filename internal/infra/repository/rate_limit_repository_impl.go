package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/beriloqueiroz/study-go-rate-limit/internal/usecase"
	redis "github.com/redis/go-redis/v9"
)

type RateLimitRepositoryImpl struct {
	Addr     string
	Password string
}

func (rr *RateLimitRepositoryImpl) FindCurrentLimiterByKey(ctx context.Context, key string) (*usecase.FindCurrentLimiterByKeyDTO, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     rr.Addr,
		Password: rr.Password,
	})
	val, err := rdb.Get(ctx, key).Bytes()

	if err != nil {
		if err.Error() == "redis: nil" {
			return &usecase.FindCurrentLimiterByKeyDTO{
				Count:    0,
				UpdateAt: time.Now(),
			}, nil
		}
		return nil, err
	}

	var output usecase.FindCurrentLimiterByKeyDTO

	json.Unmarshal(val, &output)

	return &usecase.FindCurrentLimiterByKeyDTO{
		Count:    output.Count,
		UpdateAt: output.UpdateAt,
		StartAt:  output.StartAt,
	}, nil
}

func (rr *RateLimitRepositoryImpl) Save(ctx context.Context, input *usecase.SaveInputDTO) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     rr.Addr,
		Password: rr.Password,
	})
	mr, _ := json.Marshal(&input)
	err := rdb.Set(ctx, input.Key, mr, time.Second*10)
	if err != nil {
		return err.Err()
	}
	return nil
}
