package entity

import (
	"time"
)

type limiterInfo struct {
	Key        string
	Count      int
	Limit      int
	UpdateAt   time.Time
	Expiration time.Duration
	StartAt    time.Time
}

func NewLimiterInfo(key string, count, limit int, updateAt time.Time, expiration time.Duration, startAt time.Time) *limiterInfo {
	return &limiterInfo{
		Key:        key,
		Count:      count,
		Limit:      limit,
		UpdateAt:   updateAt,
		Expiration: expiration,
		StartAt:    startAt,
	}
}

func (li *limiterInfo) isBlock() bool {
	if li.Count >= li.Limit {
		if time.Now().After(li.UpdateAt.Add(li.Expiration)) {
			li.Count = 0
			return false
		}
		return true
	}

	if li.Count == 0 {
		li.StartAt = time.Now()
	}

	li.Count = li.Count + 1

	if time.Now().Add(-time.Second).After(li.StartAt) {
		li.Count = 0
	}

	li.UpdateAt = time.Now()

	return false
}

type Limiter struct {
	IpInfo  *limiterInfo
	KeyInfo *limiterInfo
}

func NewKeyLimiter(keyInfo limiterInfo) *Limiter {
	return &Limiter{
		IpInfo:  nil,
		KeyInfo: &keyInfo,
	}
}

func NewIpLimiter(ipInfo limiterInfo) *Limiter {
	return &Limiter{
		IpInfo:  &ipInfo,
		KeyInfo: nil,
	}
}

func (l *Limiter) IsBlock() bool {
	if l.KeyInfo != nil && l.KeyInfo.Key != "" {
		return l.KeyInfo.isBlock()
	}
	if l.IpInfo != nil && l.IpInfo.Key != "" {
		return l.IpInfo.isBlock()
	}
	return false
}
