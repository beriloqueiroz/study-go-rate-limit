package usecase

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockConfigLimitRepository struct {
	mock.Mock
}

func (m *mockConfigLimitRepository) FindLimitConfigByKey(ctx context.Context, key string) (*FindLimitConfigByKeyDTO, error) {
	args := m.Called(key)
	return args.Get(0).(*FindLimitConfigByKeyDTO), args.Error(1)
}

func (m *mockConfigLimitRepository) FindLimitConfigByIp(ctx context.Context, ip string) (*FindLimitConfigByIpDTO, error) {
	args := m.Called(ip)
	return args.Get(0).(*FindLimitConfigByIpDTO), args.Error(1)
}

type mockRateLimitRepository struct {
	mock.Mock
}

func (m *mockRateLimitRepository) FindCurrentLimiterByKey(ctx context.Context, key string) (*FindCurrentLimiterByKeyDTO, error) {
	args := m.Called(key)
	return args.Get(0).(*FindCurrentLimiterByKeyDTO), args.Error(1)
}

func (m *mockRateLimitRepository) Save(ctx context.Context, input *SaveInputDTO) error {
	args := m.Called(input)
	return args.Error(0)
}
