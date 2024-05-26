package repository

import (
	"context"
	"time"

	config "github.com/beriloqueiroz/study-go-rate-limit/configs"
	"github.com/beriloqueiroz/study-go-rate-limit/internal/usecase"
)

type ConfigLimitRepositoryImpl struct {
	ConfigEnvironment *config.Conf
}

func (cr *ConfigLimitRepositoryImpl) FindLimitConfigByKey(ctx context.Context, key string) (*usecase.FindLimitConfigByKeyDTO, error) {

	var found config.ConfigDefaultApiKeysLimitPerSecond
	for _, value := range cr.ConfigEnvironment.DefaultApiKeysLimitPerSecond {
		if value.ApiKey == key {
			found = value
			break
		}
	}
	return &usecase.FindLimitConfigByKeyDTO{
		ExpirationTimeInMinutes: time.Duration(cr.ConfigEnvironment.DefaultExpirationTimeInMinutes) * time.Minute,
		LimitPerSecond:          found.LimitPerSecond,
	}, nil
}
func (cr *ConfigLimitRepositoryImpl) FindLimitConfigByIp(ctx context.Context, ip string) (*usecase.FindLimitConfigByIpDTO, error) {
	return &usecase.FindLimitConfigByIpDTO{
		ExpirationTimeInMinutes: time.Duration(cr.ConfigEnvironment.DefaultExpirationTimeInMinutes) * time.Minute,
		LimitPerSecond:          cr.ConfigEnvironment.DefaultLimitPerIPPerSecond,
	}, nil
}
