package tickets

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// durableLockRetryInterval is the bounded backoff between non-blocking lock
// attempts. Both backends poll with a non-blocking primitive so that a wait for
// the scope lock stays interruptible through the caller's context.
const durableLockRetryInterval = 2 * time.Millisecond

// durableLockBackoff waits one retry interval and reports context cancellation.
// Cancellation wins over the retry timer: a select that sees both channels ready
// picks either one, so the loop has to re-check the context after the timer fires
// instead of treating the wake-up as a plain retry.
func durableLockBackoff(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("%w: durable lock backoff requires a context", ErrInvalidCostLedger)
	}
	timer := time.NewTimer(durableLockRetryInterval)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}

// ErrDurableLockReentrant names a durable lock acquisition that would break the
// global lock order: the scope locks form a strict chain, so a goroutine that
// already holds one path may only take a strictly greater one.
var ErrDurableLockReentrant = errors.New("durable lock reentrant acquisition refused")

// durableLockOrderRegistry tracks, per goroutine, the durable lock paths it
// currently holds in acquisition order. The stack is strictly increasing by
// construction, so only the top entry has to be compared.
type durableLockOrderRegistry struct {
	mu    sync.Mutex
	stack map[uint64][]string
}

var durableLockOrder = durableLockOrderRegistry{stack: map[uint64][]string{}}

// check refuses a downgrade or a repeat: the same goroutine may only acquire a
// path beyond the deepest one it already holds.
func (registry *durableLockOrderRegistry) check(goroutine uint64, path string) error {
	if goroutine == 0 {
		return nil
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	held := registry.stack[goroutine]
	if len(held) == 0 {
		return nil
	}
	top := held[len(held)-1]
	if top >= path {
		return fmt.Errorf("%w: goroutine %d holds %q and refused %q (durable lock order must strictly increase)", ErrDurableLockReentrant, goroutine, top, path)
	}
	return nil
}

// push records a successful acquisition.
func (registry *durableLockOrderRegistry) push(goroutine uint64, path string) {
	if goroutine == 0 {
		return
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.stack[goroutine] = append(registry.stack[goroutine], path)
}

// pop drops the deepest recorded acquisition, paired with the release of that lock.
func (registry *durableLockOrderRegistry) pop(goroutine uint64) {
	if goroutine == 0 {
		return
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	held := registry.stack[goroutine]
	if len(held) == 0 {
		return
	}
	if len(held) == 1 {
		delete(registry.stack, goroutine)
		return
	}
	registry.stack[goroutine] = held[:len(held)-1]
}

// currentGoroutineID parses the identifier from the runtime stack header. It is
// the only portable way to answer "does this goroutine already hold a lock", and
// an unparsable header simply disables the order tracking for that call.
func currentGoroutineID() uint64 {
	var header [64]byte
	n := runtime.Stack(header[:], false)
	fields := strings.Fields(string(header[:n]))
	if len(fields) < 2 || fields[0] != "goroutine" {
		return 0
	}
	id, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return 0
	}
	return id
}
