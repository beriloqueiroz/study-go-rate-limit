package entity

import (
	"time"
)

type LimiterInfo struct {
	Key        string
	Count      int
	Limit      int
	UpdateAt   time.Time
	Expiration time.Duration
	StartAt    time.Time
	isBlocked  bool
}

func NewLimiterInfo(key string, count, limit int, updateAt time.Time, expiration time.Duration, startAt time.Time) *LimiterInfo {
	return &LimiterInfo{
		Key:        key,
		Count:      count,
		Limit:      limit,
		UpdateAt:   updateAt,
		Expiration: expiration,
		StartAt:    startAt,
	}
}

func (li *LimiterInfo) Process() {
	if li.Count >= li.Limit {
		if time.Now().After(li.UpdateAt.Add(li.Expiration)) {
			li.Count = 0
			li.isBlocked = false
			return
		}
		li.isBlocked = true
		return
	}

	if li.Count == 0 || li.StartAt.IsZero() {
		li.StartAt = time.Now()
	}

	li.Count = li.Count + 1

	if time.Now().Add(-time.Second).After(li.StartAt) {
		li.Count = 0
	}

	li.UpdateAt = time.Now()

	li.isBlocked = false
}

func (li *LimiterInfo) isBlock() bool {
	return li.isBlocked
}

type Limiter struct {
	IpInfo  *LimiterInfo
	KeyInfo *LimiterInfo
}

func NewKeyLimiter(keyInfo *LimiterInfo) *Limiter {
	keyInfo.Process()
	return &Limiter{
		IpInfo:  nil,
		KeyInfo: keyInfo,
	}
}

func NewIpLimiter(ipInfo *LimiterInfo) *Limiter {
	ipInfo.Process()
	return &Limiter{
		IpInfo:  ipInfo,
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
