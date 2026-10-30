package tickets

import (
	"context"
	"time"
)

// durableLockRetryInterval is the bounded backoff between non-blocking lock
// attempts. Both backends poll with a non-blocking primitive so that a wait for
// the scope lock stays interruptible through the caller's context.
const durableLockRetryInterval = 2 * time.Millisecond

// durableLockBackoff waits one retry interval and reports context cancellation.
func durableLockBackoff(ctx context.Context) error {
	if ctx == nil {
		time.Sleep(durableLockRetryInterval)
		return nil
	}
	timer := time.NewTimer(durableLockRetryInterval)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
