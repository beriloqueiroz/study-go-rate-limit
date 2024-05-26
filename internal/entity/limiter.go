package entity

import "time"

type limiterInfo struct {
	Key        string
	Count      int
	Limit      int
	UpdateAt   time.Time
	Expiration time.Duration
	IsBlock    bool
}

func NewLimiterInfo(key string, count, limit int, updateAt time.Time, expiration time.Duration, isBlock bool) *limiterInfo {
	return &limiterInfo{
		Key:        key,
		Count:      count,
		Limit:      limit,
		UpdateAt:   updateAt,
		Expiration: expiration,
		IsBlock:    isBlock,
	}
}

func (li *limiterInfo) isBlock() bool {
	if li.Count >= li.Limit {
		if li.UpdateAt.Add(li.Expiration).After(time.Now()) {
			li.Count = 0
			return false
		}
		li.IsBlock = true
		return true
	}

	li.Count = li.Count + 1

	if li.UpdateAt.After(time.Now().Add(-1 * time.Second)) {
		li.Count = 0
	}

	return false
}

type Limiter struct {
	ipInfo  *limiterInfo
	keyInfo *limiterInfo
}

func NewKeyLimiter(keyInfo limiterInfo) *Limiter {
	return &Limiter{
		ipInfo:  nil,
		keyInfo: &keyInfo,
	}
}

func NewIpLimiter(ipInfo limiterInfo) *Limiter {
	return &Limiter{
		ipInfo:  &ipInfo,
		keyInfo: nil,
	}
}

func (l *Limiter) IsBlock() bool {
	if l.keyInfo != nil && l.keyInfo.Key != "" {
		return l.keyInfo.isBlock()
	}
	if l.ipInfo != nil && l.ipInfo.Key != "" {
		return l.ipInfo.isBlock()
	}
	return false
}
