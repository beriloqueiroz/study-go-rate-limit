package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimiterInfo_IsBlock_NotBlocked_CountLessThanLimit(t *testing.T) {
	limiter := LimiterInfo{
		Key:        "test",
		Count:      5,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}
	limiter.Process()
	assert.False(t, limiter.isBlock())
}

func TestLimiterInfo_IsNotBlock_Blocked_CountEqualsLimitWithExpired(t *testing.T) {
	limiter := LimiterInfo{
		Key:        "test",
		Count:      10,
		Limit:      10,
		UpdateAt:   time.Now().Add(-25 * time.Hour),
		Expiration: 24 * time.Hour,
	}
	limiter.Process()
	assert.False(t, limiter.isBlock())
}

func TestLimiterInfo_IsBlock_Blocked_CountEqualsLimitNotExpired(t *testing.T) {
	limiter := LimiterInfo{
		Key:        "test",
		Count:      10,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}
	limiter.Process()
	assert.True(t, limiter.isBlock())
}

func TestLimiterInfo_IsBlock_Blocked_CountEqualsLimitNotExpired_InSequence(t *testing.T) {
	limiter := &LimiterInfo{
		Key:        "test",
		Count:      0,
		Limit:      5,
		UpdateAt:   time.Now().Add(-25 * time.Hour),
		Expiration: 24 * time.Hour,
	}
	limiter.Process()
	assert.False(t, limiter.isBlock())
	assert.Equal(t, 1, limiter.Count)
	limiter.Process()
	assert.False(t, limiter.isBlock())
	assert.Equal(t, 2, limiter.Count)
	limiter.Process()
	assert.False(t, limiter.isBlock())
	assert.Equal(t, 3, limiter.Count)
	limiter.Process()
	assert.False(t, limiter.isBlock())
	assert.Equal(t, 4, limiter.Count)
	limiter.Process()
	assert.False(t, limiter.isBlock())
	assert.Equal(t, 5, limiter.Count)
	limiter.Process()
	assert.True(t, limiter.isBlock())
	assert.Equal(t, 5, limiter.Count)
}

func TestLimiterInfo_IsBlock_NotBlocked_CountEqualsLimitNotExpired_InSequence(t *testing.T) {
	limiter := &LimiterInfo{
		Key:        "test",
		Count:      0,
		Limit:      5,
		UpdateAt:   time.Now().Add(-25 * time.Hour),
		Expiration: 24 * time.Hour,
	}

	limiter.Process()
	assert.False(t, limiter.isBlock())
	assert.Equal(t, 1, limiter.Count)
	limiter.Process()
	assert.False(t, limiter.isBlock())
	assert.Equal(t, 2, limiter.Count)
	limiter.Process()
	assert.False(t, limiter.isBlock())
	assert.Equal(t, 3, limiter.Count)
	limiter.Process()
	assert.False(t, limiter.isBlock())
	assert.Equal(t, 4, limiter.Count)
	time.Sleep(time.Second * 1)
	limiter.Process()
	assert.False(t, limiter.isBlock())
	assert.Equal(t, 0, limiter.Count)
}

func TestLimiterInfo_IsBlock_Blocked_CountGreaterThanLimitNotExpired(t *testing.T) {
	limiter := LimiterInfo{
		Key:        "test",
		Count:      15,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}

	limiter.Process()
	assert.True(t, limiter.isBlock())
}

func TestLimiterInfo_IsBlock_NotBlocked_CountGreaterThanLimitExpired(t *testing.T) {
	limiter := LimiterInfo{
		Key:        "test",
		Count:      15,
		Limit:      10,
		UpdateAt:   time.Now().Add(-25 * time.Hour),
		Expiration: 24 * time.Hour,
	}
	limiter.Process()
	assert.False(t, limiter.isBlock())
}

func TestNewIpLimiter(t *testing.T) {
	ipInfo := &LimiterInfo{
		Key:        "ip",
		Count:      5,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}

	limiter := NewIpLimiter(ipInfo)

	assert.Equal(t, 6, limiter.IpInfo.Count)
	assert.False(t, limiter.IsBlock())
}

func TestNewKeyLimiter(t *testing.T) {
	keyInfo := &LimiterInfo{
		Key:        "key",
		Count:      3,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}

	limiter := NewKeyLimiter(keyInfo)

	assert.Equal(t, 4, limiter.KeyInfo.Count)
	assert.False(t, limiter.IsBlock())
}

func TestLimiter_IsBlock_KeyInfoBlocked(t *testing.T) {
	keyInfo := &LimiterInfo{
		Key:        "key",
		Count:      10,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}

	limiter := NewKeyLimiter(keyInfo)

	assert.Equal(t, 10, limiter.KeyInfo.Count)
	assert.True(t, limiter.IsBlock())
}

func TestLimiter_IsBlock_IpInfoBlocked(t *testing.T) {
	ipInfo := &LimiterInfo{
		Key:        "ip",
		Count:      10,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}

	limiter := NewIpLimiter(ipInfo)

	assert.Equal(t, 10, limiter.IpInfo.Count)
	assert.True(t, limiter.IsBlock())
}

func TestLimiter_IsBlock_NotBlockedWithKey(t *testing.T) {
	ipInfo := &LimiterInfo{
		Key:        "ip",
		Count:      5,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}
	keyInfo := &LimiterInfo{
		Key:        "key",
		Count:      5,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}

	limiter := NewIpLimiter(ipInfo)

	assert.False(t, limiter.IsBlock())

	limiter = NewIpLimiter(keyInfo)

	assert.False(t, limiter.IsBlock())

}

func TestLimiter_IsBlock_NotBlockedWithIp(t *testing.T) {
	ipInfo := &LimiterInfo{
		Key:        "ip",
		Count:      5,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}
	keyInfo := &LimiterInfo{
		Key:        "",
		Count:      10,
		Limit:      10,
		UpdateAt:   time.Now(),
		Expiration: 24 * time.Hour,
	}

	limiter := NewIpLimiter(ipInfo)

	assert.False(t, limiter.IsBlock())

	limiter = NewIpLimiter(keyInfo)

	assert.False(t, limiter.IsBlock())
}

func TestLimiter_IsBlock_NoIpInfoAndKeyInfo(t *testing.T) {
	limiter := &Limiter{
		IpInfo:  nil,
		KeyInfo: nil,
	}

	assert.False(t, limiter.IsBlock())
}
