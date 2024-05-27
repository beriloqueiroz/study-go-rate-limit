package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecuteUseCase(t *testing.T) {
	mockRateLimitRepository := new(mockRateLimitRepository)
	mockRateLimitRepository.On("FindCurrentLimiterByKey", "key_a").Return(&FindCurrentLimiterByKeyDTO{}, nil)
	mockRateLimitRepository.On("Save", mock.Anything).Return(nil)
	mockConfigLimitRepository := new(mockConfigLimitRepository)
	mockConfigLimitRepository.On("FindLimitConfigByKey", "key_a").Return(&FindLimitConfigByKeyDTO{}, nil)
	mockConfigLimitRepository.On("FindLimitConfigByIp", "192.168.1.15").Return(&FindLimitConfigByIpDTO{}, nil)

	usecase := NewRateLimitUseCase(mockRateLimitRepository, mockConfigLimitRepository)

	output, err := usecase.Execute(context.Background(), RateLimitUseCaseInputDto{
		Ip:  "192.168.1.15",
		Key: "key_a",
	})
	assert.Nil(t, err)
	assert.True(t, output.Allow)
}

func TestExecuteUseCaseWhenNotBlockedByKey(t *testing.T) {
	mockRateLimitRepository := new(mockRateLimitRepository)
	mockRateLimitRepository.On("FindCurrentLimiterByKey", "key_a").Return(&FindCurrentLimiterByKeyDTO{
		Count:    0,
		UpdateAt: time.Now(),
	}, nil)
	mockRateLimitRepository.On("Save", mock.Anything).Return(nil)
	mockConfigLimitRepository := new(mockConfigLimitRepository)
	mockConfigLimitRepository.On("FindLimitConfigByKey", "key_a").Return(&FindLimitConfigByKeyDTO{
		ExpirationTimeInMinutes: 5,
		LimitPerSecond:          3,
	}, nil)
	mockConfigLimitRepository.On("FindLimitConfigByIp", "192.168.1.15").Return(&FindLimitConfigByIpDTO{}, nil)
	usecase := NewRateLimitUseCase(mockRateLimitRepository, mockConfigLimitRepository)

	output, err := usecase.Execute(context.Background(), RateLimitUseCaseInputDto{
		Ip:  "192.168.1.15",
		Key: "key_a",
	})
	assert.Nil(t, err)
	assert.True(t, output.Allow)
}

func TestExecuteUseCaseWhenBlockedByKey(t *testing.T) {
	mockRateLimitRepository := new(mockRateLimitRepository)
	mockRateLimitRepository.On("FindCurrentLimiterByKey", "key_a").Return(&FindCurrentLimiterByKeyDTO{
		Count:    3,
		UpdateAt: time.Now(),
		StartAt:  time.Now(),
	}, nil)
	mockRateLimitRepository.On("Save", mock.Anything).Return(nil)
	mockConfigLimitRepository := new(mockConfigLimitRepository)
	mockConfigLimitRepository.On("FindLimitConfigByKey", "key_a").Return(&FindLimitConfigByKeyDTO{
		ExpirationTimeInMinutes: time.Minute * 5,
		LimitPerSecond:          3,
	}, nil)
	mockConfigLimitRepository.On("FindLimitConfigByIp", "192.168.1.15").Return(&FindLimitConfigByIpDTO{}, nil)
	usecase := NewRateLimitUseCase(mockRateLimitRepository, mockConfigLimitRepository)

	output, err := usecase.Execute(context.Background(), RateLimitUseCaseInputDto{
		Ip:  "192.168.1.15",
		Key: "key_a",
	})
	assert.Nil(t, err)
	assert.False(t, output.Allow)
}

func TestExecuteUseCaseWhenNotBlockedByIp(t *testing.T) {
	mockRateLimitRepository := new(mockRateLimitRepository)
	mockRateLimitRepository.On("FindCurrentLimiterByKey", "192.168.1.15").Return(&FindCurrentLimiterByKeyDTO{
		Count:    0,
		UpdateAt: time.Now(),
		StartAt:  time.Now(),
	}, nil)
	mockRateLimitRepository.On("Save", mock.Anything).Return(nil)
	mockConfigLimitRepository := new(mockConfigLimitRepository)
	mockConfigLimitRepository.On("FindLimitConfigByKey", "key_a").Return(&FindLimitConfigByKeyDTO{}, nil)
	mockConfigLimitRepository.On("FindLimitConfigByIp", "192.168.1.15").Return(&FindLimitConfigByIpDTO{
		ExpirationTimeInMinutes: time.Minute * 5,
		LimitPerSecond:          3,
	}, nil)
	usecase := NewRateLimitUseCase(mockRateLimitRepository, mockConfigLimitRepository)

	output, err := usecase.Execute(context.Background(), RateLimitUseCaseInputDto{
		Ip: "192.168.1.15",
	})
	assert.Nil(t, err)
	assert.True(t, output.Allow)
}

func TestExecuteUseCaseWhenBlockedByIp(t *testing.T) {
	mockRateLimitRepository := new(mockRateLimitRepository)
	mockRateLimitRepository.On("FindCurrentLimiterByKey", "192.168.1.15").Return(&FindCurrentLimiterByKeyDTO{
		Count:    3,
		UpdateAt: time.Now(),
		StartAt:  time.Now(),
	}, nil)
	mockRateLimitRepository.On("Save", mock.Anything).Return(nil)
	mockConfigLimitRepository := new(mockConfigLimitRepository)
	mockConfigLimitRepository.On("FindLimitConfigByKey", "key_a").Return(&FindLimitConfigByKeyDTO{}, nil)
	mockConfigLimitRepository.On("FindLimitConfigByIp", "192.168.1.15").Return(&FindLimitConfigByIpDTO{
		ExpirationTimeInMinutes: time.Minute * 5,
		LimitPerSecond:          3,
	}, nil)
	usecase := NewRateLimitUseCase(mockRateLimitRepository, mockConfigLimitRepository)

	output, err := usecase.Execute(context.Background(), RateLimitUseCaseInputDto{
		Ip: "192.168.1.15",
	})
	assert.Nil(t, err)
	assert.False(t, output.Allow)
}
