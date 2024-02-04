package redis_lock

import (
	"time"
)

type Retry interface {
	Next() (time.Duration, bool)
}

type TestRetry struct {
	duration time.Duration
	maxCnt   int
	cnt      int
}

func (t *TestRetry) Next() (time.Duration, bool) {
	if t.cnt >= t.maxCnt {
		return 0, false
	}
	t.cnt++
	return t.duration, true
}
