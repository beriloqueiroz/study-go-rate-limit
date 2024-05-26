package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimiterInfo_IsBlock_NotBlocked_CountLessThanLimit(t *testing.T) {
	limiter := limiterInfo{
		Key:        "test",
		Count:      5,
		Limit:      10,
		UpdateAt:   time.Now().Add(-1 * time.Hour),
		Expiration: 24 * time.Hour,
	}

	assert.False(t, limiter.isBlock())
}

func TestLimiterInfo_IsBlock_NotBlocked_CountEqualsLimitExpired(t *testing.T) {
	limiter := limiterInfo{
		Key:        "test",
		Count:      10,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}

	assert.False(t, limiter.isBlock())
}

func TestLimiterInfo_IsBlock_Blocked_CountEqualsLimitNotExpired(t *testing.T) {
	limiter := limiterInfo{
		Key:        "test",
		Count:      10,
		Limit:      10,
		UpdateAt:   time.Now().Add(-25 * time.Hour),
		Expiration: 24 * time.Hour,
	}

	assert.True(t, limiter.isBlock())
}

func TestLimiterInfo_IsBlock_Blocked_CountGreaterThanLimitNotExpired(t *testing.T) {
	limiter := limiterInfo{
		Key:        "test",
		Count:      15,
		Limit:      10,
		UpdateAt:   time.Now().Add(-25 * time.Hour),
		Expiration: 24 * time.Hour,
	}

	assert.True(t, limiter.isBlock())
}

func TestLimiterInfo_IsBlock_NotBlocked_CountGreaterThanLimitExpired(t *testing.T) {
	limiter := limiterInfo{
		Key:        "test",
		Count:      15,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}

	assert.False(t, limiter.isBlock())
}

func TestNewLimiter(t *testing.T) {
	ipInfo := limiterInfo{
		Key:        "ip",
		Count:      5,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}
	keyInfo := limiterInfo{
		Key:        "key",
		Count:      3,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}

	limiter := NewLimiter(ipInfo, keyInfo)

	assert.Equal(t, ipInfo, *limiter.ipInfo)
	assert.Equal(t, keyInfo, *limiter.keyInfo)
}

func TestLimiter_IsBlock_KeyInfoBlocked(t *testing.T) {
	keyInfo := &limiterInfo{
		Key:        "key",
		Count:      10,
		Limit:      10,
		UpdateAt:   time.Now().Add(-25 * time.Hour),
		Expiration: 24 * time.Hour,
	}

	limiter := &Limiter{
		ipInfo:  nil,
		keyInfo: keyInfo,
	}

	assert.True(t, limiter.IsBlock())
}

func TestLimiter_IsBlock_KeyInfoBlockedWithIp(t *testing.T) {
	ipInfo := &limiterInfo{
		Key:        "ip",
		Count:      5,
		Limit:      10,
		UpdateAt:   time.Now().Add(-1 * time.Hour),
		Expiration: 24 * time.Hour,
	}
	keyInfo := &limiterInfo{
		Key:        "key",
		Count:      10,
		Limit:      10,
		UpdateAt:   time.Now().Add(-25 * time.Hour),
		Expiration: 24 * time.Hour,
	}

	limiter := &Limiter{
		ipInfo:  ipInfo,
		keyInfo: keyInfo,
	}

	assert.True(t, limiter.IsBlock())
}

func TestLimiter_IsBlock_IpInfoBlocked(t *testing.T) {
	ipInfo := &limiterInfo{
		Key:        "ip",
		Count:      10,
		Limit:      10,
		UpdateAt:   time.Now().Add(-25 * time.Hour),
		Expiration: 24 * time.Hour,
	}

	limiter := &Limiter{
		ipInfo:  ipInfo,
		keyInfo: nil,
	}

	assert.True(t, limiter.IsBlock())
}

func TestLimiter_IsBlock_NotBlockedWithKey(t *testing.T) {
	ipInfo := &limiterInfo{
		Key:        "ip",
		Count:      5,
		Limit:      10,
		UpdateAt:   time.Now().Add(-1 * time.Hour),
		Expiration: 24 * time.Hour,
	}
	keyInfo := &limiterInfo{
		Key:        "key",
		Count:      5,
		Limit:      10,
		UpdateAt:   time.Now().Add(-1 * time.Hour),
		Expiration: 24 * time.Hour,
	}

	limiter := &Limiter{
		ipInfo:  ipInfo,
		keyInfo: keyInfo,
	}

	assert.False(t, limiter.IsBlock())
}

func TestLimiter_IsBlock_NotBlockedWithIp(t *testing.T) {
	ipInfo := &limiterInfo{
		Key:        "ip",
		Count:      5,
		Limit:      10,
		UpdateAt:   time.Now().Add(-1 * time.Hour),
		Expiration: 24 * time.Hour,
	}
	keyInfo := &limiterInfo{
		Key:        "",
		Count:      10,
		Limit:      10,
		UpdateAt:   time.Now().Add(-25 * time.Hour),
		Expiration: 24 * time.Hour,
	}

	limiter := &Limiter{
		ipInfo:  ipInfo,
		keyInfo: keyInfo,
	}

	assert.False(t, limiter.IsBlock())
}

func TestLimiter_IsBlock_NoIpInfoAndKeyInfo(t *testing.T) {
	limiter := &Limiter{
		ipInfo:  nil,
		keyInfo: nil,
	}

	assert.False(t, limiter.IsBlock())
}
