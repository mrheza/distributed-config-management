package service

import (
	"context"
	"math/rand"
	"time"
)

type retryState struct {
	count      int
	lastTarget string
}

func (r *retryState) reset() {
	r.count = 0
	r.lastTarget = ""
}

func (r *retryState) next(target string) int {
	if target == "" {
		target = "remote"
	}

	if r.lastTarget != "" && r.lastTarget != target {
		r.count = 0
	}
	r.lastTarget = target
	r.count++

	return r.count
}

func sleepWithContext(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		d = time.Millisecond
	}
	select {
	case <-ctx.Done():
		return false
	case <-time.After(d):
		return true
	}
}

func calculateBackoff(retryCount, maxBackoffSecs int) time.Duration {
	if retryCount < 1 {
		retryCount = 1
	}

	maxBackoff := time.Duration(maxBackoffSecs) * time.Second
	if maxBackoff <= 0 {
		maxBackoff = time.Second
	}

	backoff := time.Second
	for i := 1; i < retryCount; i++ {
		if backoff >= maxBackoff {
			return maxBackoff
		}
		backoff *= 2
	}

	if backoff > maxBackoff {
		return maxBackoff
	}
	return backoff
}

func applyJitter(base time.Duration, jitterPercent int) time.Duration {
	return applyJitterWithFloat64(base, jitterPercent, rand.Float64)
}

func (s *agentService) applyJitter(base time.Duration, jitterPercent int) time.Duration {
	if s.rng == nil {
		return applyJitter(base, jitterPercent)
	}
	return applyJitterWithFloat64(base, jitterPercent, s.rng.Float64)
}

func applyJitterWithFloat64(base time.Duration, jitterPercent int, nextFloat64 func() float64) time.Duration {
	if base <= 0 || jitterPercent <= 0 {
		return base
	}

	if jitterPercent > 90 {
		jitterPercent = 90
	}

	delta := float64(base) * float64(jitterPercent) / 100.0
	min := float64(base) - delta
	max := float64(base) + delta
	jittered := min + nextFloat64()*(max-min)
	if jittered < 0 {
		return 0
	}

	return time.Duration(jittered)
}
