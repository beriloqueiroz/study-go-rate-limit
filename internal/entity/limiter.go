package entity

import "time"

type limiterInfo struct {
	Key        string
	Count      int
	Limit      int
	UpdateAt   time.Time
	Expiration time.Duration
}

func (li *limiterInfo) isBlock() bool {
	if li.Count >= li.Limit {
		if li.UpdateAt.Add(li.Expiration).After(time.Now()) {
			li.Count = 0
			return false
		}
		return true
	}
	return false
}

type Limiter struct {
	ipInfo  *limiterInfo
	keyInfo *limiterInfo
}

func NewLimiter(ipInfo, keyInfo limiterInfo) *Limiter {
	return &Limiter{
		ipInfo:  &ipInfo,
		keyInfo: &keyInfo,
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
