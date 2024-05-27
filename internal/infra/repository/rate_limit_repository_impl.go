package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/beriloqueiroz/study-go-rate-limit/internal/usecase"
	redis "github.com/redis/go-redis/v9"
)

type RateLimitRepositoryImpl struct {
}

func (rr *RateLimitRepositoryImpl) FindCurrentLimiterByKey(ctx context.Context, key string) (*usecase.FindCurrentLimiterByKeyDTO, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "my-password",
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

	fmt.Println(output, key)

	return &usecase.FindCurrentLimiterByKeyDTO{
		Count:    output.Count,
		UpdateAt: output.UpdateAt,
		StartAt:  output.StartAt,
	}, nil
}

func (rr *RateLimitRepositoryImpl) Save(ctx context.Context, input *usecase.SaveInputDTO) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "my-password",
	})
	mr, _ := json.Marshal(&input)
	err := rdb.Set(ctx, input.Key, mr, time.Second*10)
	if err != nil {
		return err.Err()
	}
	return nil
}
